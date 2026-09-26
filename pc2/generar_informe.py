"""Build the integrated PC2 report PDF from the team's PC1 PDF and PC2 Markdown."""

from __future__ import annotations

import html
import tempfile
import textwrap
from pathlib import Path

import fitz
import mistune
from bs4 import BeautifulSoup, NavigableString, Tag
from PIL import Image as PILImage, ImageDraw, ImageFont
from reportlab.lib import colors
from reportlab.lib.enums import TA_CENTER, TA_LEFT
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle, getSampleStyleSheet
from reportlab.lib.units import mm
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import (
    HRFlowable,
    Image,
    ListFlowable,
    ListItem,
    PageBreak,
    Paragraph,
    Preformatted,
    SimpleDocTemplate,
    Spacer,
    Table,
    TableStyle,
)


ROOT = Path(__file__).resolve().parents[1]
PC2 = Path(__file__).resolve().parent
OLD_PDF = ROOT / "CC65-PC1-202620.pdf"
PC2_MD = PC2 / "INFORME_PC2.md"
FINAL_PDF = PC2 / "Informe_CC65-PC2-202620_Equipo.pdf"
TMP_DIR = Path(tempfile.gettempdir())
APPENDIX_PDF = TMP_DIR / "pc2-appendix.pdf"
FRONT_PDF = TMP_DIR / "pc2-front.pdf"
EVIDENCE_PNG = PC2 / "evidencia_ejecucion.png"


def register_fonts() -> tuple[str, str, str]:
    regular = Path("C:/Windows/Fonts/arial.ttf")
    bold = Path("C:/Windows/Fonts/arialbd.ttf")
    italic = Path("C:/Windows/Fonts/ariali.ttf")
    if regular.exists() and bold.exists() and italic.exists():
        pdfmetrics.registerFont(TTFont("Arial", str(regular)))
        pdfmetrics.registerFont(TTFont("Arial-Bold", str(bold)))
        pdfmetrics.registerFont(TTFont("Arial-Italic", str(italic)))
        return "Arial", "Arial-Bold", "Arial-Italic"
    return "Helvetica", "Helvetica-Bold", "Helvetica-Oblique"


BODY_FONT, BOLD_FONT, ITALIC_FONT = register_fonts()
NAVY = colors.HexColor("#16324F")
BLUE = colors.HexColor("#2E6F95")
PALE = colors.HexColor("#EAF1F6")
INK = colors.HexColor("#202832")


def make_styles():
    styles = getSampleStyleSheet()
    styles.add(ParagraphStyle(name="BodyES", parent=styles["BodyText"], fontName=BODY_FONT, fontSize=9.2, leading=13.2, textColor=INK, spaceAfter=6))
    styles.add(ParagraphStyle(name="H1ES", parent=styles["Heading1"], fontName=BOLD_FONT, fontSize=17, leading=21, textColor=NAVY, spaceBefore=8, spaceAfter=10, keepWithNext=True))
    styles.add(ParagraphStyle(name="H2ES", parent=styles["Heading2"], fontName=BOLD_FONT, fontSize=12.2, leading=15, textColor=BLUE, spaceBefore=7, spaceAfter=6, keepWithNext=True))
    styles.add(ParagraphStyle(name="H3ES", parent=styles["Heading3"], fontName=BOLD_FONT, fontSize=10.5, leading=13, textColor=NAVY, spaceBefore=5, spaceAfter=5, keepWithNext=True))
    styles.add(ParagraphStyle(name="SmallES", parent=styles["BodyText"], fontName=BODY_FONT, fontSize=8.1, leading=10.3, textColor=INK, spaceAfter=4))
    styles.add(ParagraphStyle(name="CodeES", fontName="Courier", fontSize=7.8, leading=10.2, leftIndent=8, rightIndent=8, borderColor=colors.HexColor("#D7E0E8"), borderWidth=0.5, borderPadding=6, backColor=colors.HexColor("#F5F7F9"), spaceBefore=4, spaceAfter=8))
    styles.add(ParagraphStyle(name="CoverES", fontName=BOLD_FONT, fontSize=21, leading=28, alignment=TA_CENTER, textColor=NAVY, spaceAfter=12))
    styles.add(ParagraphStyle(name="CoverSubES", fontName=BODY_FONT, fontSize=12, leading=17, alignment=TA_CENTER, textColor=BLUE, spaceAfter=7))
    styles.add(ParagraphStyle(name="TOCES", fontName=BODY_FONT, fontSize=10, leading=16, leftIndent=4, spaceAfter=3))
    return styles


STYLES = make_styles()


def page_footer(canvas, doc):
    canvas.saveState()
    canvas.setStrokeColor(colors.HexColor("#D7E0E8"))
    canvas.line(20 * mm, 15 * mm, A4[0] - 20 * mm, 15 * mm)
    canvas.setFont(BODY_FONT, 8)
    canvas.setFillColor(colors.HexColor("#667788"))
    page = canvas.getPageNumber() + doc.page_offset
    canvas.drawString(20 * mm, 10 * mm, "CC65 · Informe PC2 · PaySim1")
    canvas.drawRightString(A4[0] - 20 * mm, 10 * mm, f"Página {page}")
    canvas.restoreState()


def sanitize_inline(source: str) -> str:
    source = source.replace("<code>", f'<font name="Courier">').replace("</code>", "</font>")
    source = source.replace("<del>", "<strike>").replace("</del>", "</strike>")
    return source


def to_flowables(markdown: str):
    renderer = mistune.create_markdown(renderer="html", plugins=["table", "strikethrough"])
    soup = BeautifulSoup(renderer(markdown), "html.parser")
    story = []
    for item in soup.contents:
        if isinstance(item, NavigableString) or not isinstance(item, Tag):
            continue
        tag = item.name.lower()
        if tag in {"h1", "h2", "h3", "h4"}:
            key = {"h1": "H1ES", "h2": "H2ES", "h3": "H3ES", "h4": "H3ES"}[tag]
            story.append(Paragraph(sanitize_inline(item.decode_contents()), STYLES[key]))
        elif tag == "p":
            image_tag = item.find("img")
            if image_tag is not None:
                path = PC2 / Path(image_tag.get("src", "")).name
                if not path.exists():
                    raise FileNotFoundError(f"report image not found: {path}")
                pic = Image(str(path), width=A4[0] - 42 * mm, height=100 * mm, kind="proportional")
                pic.hAlign = "CENTER"
                story.extend([pic, Paragraph(html.escape(image_tag.get("alt", "Evidencia de ejecución")), STYLES["SmallES"]), Spacer(1, 6)])
            else:
                story.append(Paragraph(sanitize_inline(item.decode_contents()), STYLES["BodyES"]))
        elif tag in {"ul", "ol"}:
            entries = []
            for li in item.find_all("li", recursive=False):
                entries.append(ListItem(Paragraph(sanitize_inline(li.decode_contents()), STYLES["BodyES"]), leftIndent=10))
            story.append(ListFlowable(entries, bulletType="1" if tag == "ol" else "bullet", start="1", leftIndent=17, bulletFontName=BODY_FONT, bulletFontSize=8, spaceAfter=4))
        elif tag == "pre":
            text = item.get_text().rstrip()
            story.append(Preformatted(text, STYLES["CodeES"], maxLineLength=100))
        elif tag == "table":
            rows = []
            for row_number, tr in enumerate(item.find_all("tr")):
                cells = []
                for td in tr.find_all(["th", "td"], recursive=False):
                    content = sanitize_inline(td.decode_contents()).replace("<p>", "").replace("</p>", "<br/>")
                    if row_number == 0:
                        header_style = ParagraphStyle(name="TableHead", parent=STYLES["SmallES"], fontName=BOLD_FONT, textColor=colors.white)
                        cells.append(Paragraph(content, header_style))
                    else:
                        cells.append(Paragraph(content, STYLES["SmallES"]))
                if cells:
                    rows.append(cells)
            if rows:
                ncols = max(map(len, rows))
                widths = [(A4[0] - 42 * mm) / ncols] * ncols
                table = Table(rows, colWidths=widths, repeatRows=1, hAlign="LEFT")
                table.setStyle(TableStyle([
                    ("BACKGROUND", (0, 0), (-1, 0), NAVY),
                    ("TEXTCOLOR", (0, 0), (-1, 0), colors.white),
                    ("FONTNAME", (0, 0), (-1, 0), BOLD_FONT),
                    ("BACKGROUND", (0, 1), (-1, -1), colors.white),
                    ("ROWBACKGROUNDS", (0, 1), (-1, -1), [colors.white, PALE]),
                    ("GRID", (0, 0), (-1, -1), 0.35, colors.HexColor("#C7D2DC")),
                    ("VALIGN", (0, 0), (-1, -1), "TOP"),
                    ("LEFTPADDING", (0, 0), (-1, -1), 5),
                    ("RIGHTPADDING", (0, 0), (-1, -1), 5),
                    ("TOPPADDING", (0, 0), (-1, -1), 5),
                    ("BOTTOMPADDING", (0, 0), (-1, -1), 3),
                ]))
                story.extend([table, Spacer(1, 6)])
        elif tag == "hr":
            story.append(HRFlowable(width="100%", thickness=0.7, color=BLUE, spaceBefore=5, spaceAfter=8))
        elif tag == "blockquote":
            story.append(Paragraph(sanitize_inline(item.get_text(" ", strip=True)), STYLES["BodyES"]))
        elif tag == "img":
            path = PC2 / Path(item.get("src", "")).name
            if not path.exists():
                raise FileNotFoundError(f"report image not found: {path}")
            pic = Image(str(path), width=A4[0] - 42 * mm, height=100 * mm, kind="proportional")
            pic.hAlign = "CENTER"
            story.extend([pic, Spacer(1, 7)])
    return story


def make_appendix():
    markdown = PC2_MD.read_text(encoding="utf-8")
    story = to_flowables(markdown)
    # The PC2 section begins after the two front pages and 17 retained PC1 pages.
    doc = SimpleDocTemplate(str(APPENDIX_PDF), pagesize=A4, rightMargin=21 * mm, leftMargin=21 * mm, topMargin=19 * mm, bottomMargin=22 * mm, title="CC65 PC2 - PaySim1")
    doc.page_offset = 19
    doc.build(story, onFirstPage=page_footer, onLaterPages=page_footer)


def make_front(appendix_pages: int):
    entries = [
        "1. Resumen del trabajo ............................................................ 3",
        "2. Objetivos del trabajo ............................................................ 3",
        "3. Investigación bibliográfica .................................................... 4",
        "4. Análisis del caso de uso ...................................................... 13",
        "5. Limpieza y preprocesamiento de PaySim1 ............................ 15",
        "6. Repositorio e historial Git .................................................. 19",
        "7. PC2: modelo, implementación y evaluación ....................... 20",
        f"8. Referencias ........................................................................ {20 + appendix_pages}",
    ]
    story = [Spacer(1, 23 * mm), Paragraph("UNIVERSIDAD PERUANA DE CIENCIAS APLICADAS", STYLES["CoverSubES"]), Paragraph("FACULTAD DE INGENIERÍA · CIENCIAS DE LA COMPUTACIÓN", STYLES["CoverSubES"]), Spacer(1, 9 * mm), HRFlowable(width="62%", thickness=1.5, color=BLUE, hAlign="CENTER"), Spacer(1, 10 * mm), Paragraph("INFORME PC2", STYLES["CoverES"]), Paragraph("Detección concurrente de fraude en transacciones PaySim1", STYLES["CoverSubES"]), Spacer(1, 3 * mm), Paragraph("Programación Concurrente y Distribuida · 1ACC0065", STYLES["CoverSubES"]), Paragraph("Docente: Herminio Paucar Curasma · NRC: 8646", STYLES["CoverSubES"]), Spacer(1, 9 * mm), Paragraph("INTEGRANTES", STYLES["H3ES"])]
    members = [
        "Avalos Sánchez, César Gabriel (u202310307)",
        "Rojas Cuadros, Fabian Marcelo (u202218498)",
        "Ballon Villar, Diego Eduardo (u201520327)",
        "Cuadrado Jimenez, Sebastian Antonio (u202312882)",
    ]
    story.extend(Paragraph(name, ParagraphStyle("Member" + str(i), parent=STYLES["CoverSubES"], fontSize=10, leading=14, spaceAfter=2)) for i, name in enumerate(members))
    story.extend([Spacer(1, 7 * mm), Paragraph("Ciclo académico 202620", STYLES["CoverSubES"]), PageBreak(), Paragraph("Índice", STYLES["H1ES"])])
    story.extend(Paragraph(html.escape(line), STYLES["TOCES"]) for line in entries)
    doc = SimpleDocTemplate(str(FRONT_PDF), pagesize=A4, rightMargin=24 * mm, leftMargin=24 * mm, topMargin=22 * mm, bottomMargin=21 * mm, title="CC65 PC2 - Índice")
    doc.page_offset = 0
    doc.build(story, onFirstPage=page_footer, onLaterPages=page_footer)
    front = fitz.open(str(FRONT_PDF))
    if len(front) != 2:
        raise RuntimeError(f"front section should be 2 pages; got {len(front)}")
    front.close()


def create_evidence_image():
    bench = (PC2 / "registro_rendimiento.txt").read_text(encoding="utf-16").splitlines()
    evaluation = (PC2 / "registro_evaluacion.txt").read_text(encoding="utf-16").splitlines()
    summary = [line for line in bench if "trimmed_mean_10pct=" in line]
    resource = next((line for line in bench if line.startswith("Wrote raw observations")), "")
    dataset_lines = [line for line in bench if line.startswith(("Dataset:", "Training frauds:", "Benchmark:"))]
    metrics_lines = [line for line in evaluation if line.startswith("workers=")]
    lines = [
        "PS  PC2  |  stdout real capturado",
        "> .\\pc2.exe -input ..\\paysim.csv -mode bench -repeats 30 -bench-epochs 1 -workers 1,2,4,8",
        *dataset_lines,
        *summary,
        resource,
        "> .\\pc2.exe -input ..\\paysim.csv -mode evaluate -epochs 5",
        *metrics_lines,
    ]
    lines = [line for line in lines if line]
    font_path = Path("C:/Windows/Fonts/consola.ttf")
    bold_path = Path("C:/Windows/Fonts/consolab.ttf")
    font = ImageFont.truetype(str(font_path), 23) if font_path.exists() else ImageFont.load_default()
    bold = ImageFont.truetype(str(bold_path), 24) if bold_path.exists() else font
    lines = [wrapped for line in lines for wrapped in textwrap.wrap(line, width=75, subsequent_indent="    ")]
    width = 1280
    line_height = 37
    height = 66 + len(lines) * line_height + 22
    img = PILImage.new("RGB", (width, height), "#101820")
    draw = ImageDraw.Draw(img)
    draw.rounded_rectangle((18, 14, width - 18, height - 14), radius=12, fill="#17232e", outline="#355066", width=2)
    draw.ellipse((36, 32, 50, 46), fill="#ff6b6b")
    draw.ellipse((60, 32, 74, 46), fill="#ffd166")
    draw.ellipse((84, 32, 98, 46), fill="#06d6a0")
    y = 62
    for i, line in enumerate(lines):
        draw.text((38, y), line, font=bold if i == 0 else font, fill="#9fe7bf" if line.startswith("workers=") and "elapsed=" in line else "#e6edf3")
        y += line_height
    img.save(EVIDENCE_PNG, optimize=True)


def assemble_pdf(appendix_pages: int):
    make_front(appendix_pages)
    old = fitz.open(str(OLD_PDF))
    front = fitz.open(str(FRONT_PDF))
    appendix = fitz.open(str(APPENDIX_PDF))
    final = fitz.open()
    final.insert_pdf(front)
    # Keep all PC1 report pages except its cover, old index, and references page.
    final.insert_pdf(old, from_page=2, to_page=len(old) - 2)
    final.insert_pdf(appendix)
    # Move existing PC1 references to the end of the integrated report.
    final.insert_pdf(old, from_page=len(old) - 1, to_page=len(old) - 1)
    # The new contents use section number 8; correct the retained PC1 page heading.
    ref_page = final[len(final) - 1]
    for block in ref_page.get_text("dict")["blocks"]:
        for line in block.get("lines", []):
            for span in line.get("spans", []):
                if "Referencias." in span.get("text", ""):
                    box = fitz.Rect(span["bbox"])
                    ref_page.add_redact_annot(fitz.Rect(box.x0 - 1, box.y0 - 1, box.x1 + 1, box.y1 + 1), fill=(1, 1, 1))
                    ref_page.apply_redactions()
                    ref_page.insert_text((box.x0, box.y1 - 2), "8. Referencias.", fontsize=span["size"], fontname="helv", color=(0, 0, 0))
                    break
    final.set_metadata({"title": "CC65-PC2-202620 - Detección concurrente de fraude en PaySim1", "author": "Equipo CC65", "subject": "Informe integrado PC1 y PC2"})
    PC2.mkdir(parents=True, exist_ok=True)
    final.save(str(FINAL_PDF), garbage=4, deflate=True)
    final.close(); appendix.close(); front.close(); old.close()


def main():
    create_evidence_image()
    global TMP_DIR, APPENDIX_PDF, FRONT_PDF
    with tempfile.TemporaryDirectory(prefix="cc65-pc2-report-") as scratch:
        TMP_DIR = Path(scratch)
        APPENDIX_PDF = TMP_DIR / "pc2-appendix.pdf"
        FRONT_PDF = TMP_DIR / "pc2-front.pdf"
        make_appendix()
        appendix = fitz.open(str(APPENDIX_PDF))
        appendix_pages = len(appendix)
        appendix.close()
        assemble_pdf(appendix_pages)
        print(f"Created {FINAL_PDF} with {appendix_pages} PC2 pages")


if __name__ == "__main__":
    main()
