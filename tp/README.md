# TP de detección concurrente de fraude

El TP añade a PC2 verificación formal ejecutada, análisis IA y correcciones de validación. La rama de trabajo es `feature/tp`; su integración en `main` continúa pendiente.

## Archivos de entrega

- [Informe IA con el prompt estructurado incluido como anexo](INFORME_IA.md).
- [Verificación formal](VERIFICACION_SPIN.md), [modelo](modelo_tp.pml), [resultados](evidencias/spin_resumen.csv) y logs completos en `evidencias/`.
- Word integrado local en la raíz: `CC65-TP-202620-[Código de alumno].docx`. Cada alumno debe renombrar su copia y validar su análisis personal.
- Plantilla Word local para el coordinador: `CC65-Participación-202620.docx`. Completar solo participación real.

El Word de PC2 se conserva sin cambios. Word, PDF y revisión visual permanecen fuera de Git. El informe IA y las evidencias TP sí deben publicarse posteriormente.

## Reproducir Spin

Desde la raíz, con Docker Desktop funcionando:

```powershell
docker build -t cc65-tp-spin:local tp
$tpMount = 'type=bind,source=' + (Resolve-Path tp).Path + ',target=/tp'
docker run --rm --network none --mount $tpMount cc65-tp-spin:local
```

El contenedor solo monta `tp/` y la ejecución no usa red. La construcción obtiene paquetes oficiales Debian y fija el digest de la imagen base. Las versiones ejecutadas están en `evidencias/spin_entorno.txt`; futuras actualizaciones de paquetes pueden cambiar sus versiones. El script falla si un modelo correcto tiene errores o búsqueda incompleta, o si un control defectuoso no produce un contraejemplo.

Con Spin, GCC y Python instalados en Linux:

```bash
python3 tp/verificar_spin.py
```

## Reproducir Go

Desde `pc2/`, con Go y GCC disponibles:

```powershell
$env:CGO_ENABLED = '1'
$env:CC = 'gcc'
go test -race -count=1 -v ./...
go vet ./...
$env:TP_FULL_DATASET = '..\paysim.csv'
go test -run TestFullPaySim -count=1 -v ./...
Remove-Item Env:TP_FULL_DATASET
```

La integración completa es optativa en las pruebas ordinarias porque lee 6,362,620 filas. Se ejecutó por separado sin detector de carreras. La equivalencia y el detector sí cubren conjuntos pequeños, particiones desiguales y más workers que filas.

Windows bloqueó algunos ejecutables, incluyendo Spin precompilado y ciertas invocaciones directas Go. Se usó Docker para Spin y `go test` con GCC portátil para Go. Los intentos bloqueados se conservan identificados y no se cuentan como éxitos.

## Historial y publicación posterior

Commit público analizado: `0a3495740ae7c5827a2ad8095ef33677ba82519b`. Consultar `git log --oneline origin/feature/pc2..feature/tp` para obtener el historial local final. No inventar autorías ni participación.

Después de revisar `feature/tp`, integrar mediante el flujo de PR del equipo hasta `main` antes de la fecha límite. Guardar enlace al commit final y captura real de la rama principal. La consigna llama `master` a esa rama; este repositorio usa `main`.

La captura de PC2 conservada en el Word es histórica y no acredita la publicación de TP. El video y su enlace deben generarse por el equipo.
