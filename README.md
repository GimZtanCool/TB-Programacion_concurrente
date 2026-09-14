# Detección de Fraude en Transacciones Financieras Masivas (PaySim1)
### CC65 - Programación Concurrente y Distribuida | Trabajo Parcial 2026-20

Repositorio del proyecto del curso **Programación Concurrente y Distribuida (CC65)**, enfocado en el procesamiento masivo de datos, análisis exploratorio (EDA) y preparación de un conjunto de datos a gran escala para el entrenamiento e inferencia de modelos de Machine Learning mediante arquitecturas concurrentes y paralelas (Go, Goroutines, Worker Pools y verificación formal en Promela/Spin).

---

## 📌 1. Contexto y Objetivos del Proyecto

El fraude en transacciones electrónicas representa pérdidas multimillonarias anuales a nivel global y vulnera la confianza en la infraestructura financiera digital. Los objetivos del proyecto son:
1. **Analizar y preprocesar** un conjunto de datos masivo (>6.3 millones de transacciones con variables financieras, contables y temporales) garantizando consistencia, balanceo y escalabilidad.
2. **Implementar modelos de Machine Learning** bajo esquemas secuenciales y concurrentes en lenguaje **Go** (sin librerías de terceros para la concurrencia), explotando mecanismos nativos como `goroutines`, canales (`chan`), `sync.WaitGroup` y `sync.Mutex`.
3. **Evaluar el rendimiento computacional** mediante métricas de *Speedup*, escalabilidad y medias recortadas, modelando la sincronización en **Promela/Spin** para verificar formalmente la ausencia de condiciones de carrera y *deadlocks*.

### Alineación con los ODS (Objetivos de Desarrollo Sostenible)
* **ODS 8: Trabajo decente y crecimiento económico:** Fortalecimiento de la seguridad en sistemas de pago digitales e inclusión financiera confiable.
* **ODS 9: Industria, innovación e infraestructura:** Construcción de infraestructura de software robusta, concurrente y de alta disponibilidad para el procesamiento de transacciones en tiempo real.
* **ODS 16: Paz, justicia e instituciones sólidas:** Reducción sustancial del flujo financiero ilícito y actividades fraudulentas en el comercio electrónico.

---

## 📊 2. Dataset Seleccionado: Synthetic Financial Datasets For Fraud Detection (PaySim1)

El dataset seleccionado corresponde a **PaySim1**, un simulador financiero basado en datos transaccionales reales de servicios móviles de dinero:
* **Fuente:** [PaySim1 en Kaggle](https://www.kaggle.com/datasets/ealaxi/paysim1)
* **Volumen:** **6,362,620 transacciones** y 11 variables (cumpliendo con creces el requisito de >1M de registros).
* **Variable objetivo:** `isFraud` (binaria: 0 = transacción legítima, 1 = fraude).
  * **Distribución:** 6,354,407 transacciones legítimas (99.8709%) vs. 8,213 fraudulentas (0.1291%), reflejando el severo desbalance de clases inherente al fraude bancario real.
  * **Hallazgo clave de fraude:** El fraude se concentra **exclusivamente** en los tipos `TRANSFER` (4,097 casos) y `CASH_OUT` (4,116 casos). En `CASH_IN`, `DEBIT` y `PAYMENT` la tasa de fraude es exactamente 0.

---

## 🛠️ 3. Pipeline de Limpieza y Preprocesamiento (`limpieza_paysim_fraud.ipynb`)

El procedimiento implementado en el notebook sigue un flujo estricto y reproducible:

1. **Diagnóstico de Calidad e Integridad:**
   * Verificación de valores faltantes: **0 nulos** en todas las variables.
   * Verificación de duplicados: **0 registros duplicados**.
2. **Tratamiento de Identificadores:**
   * Depuración de identificadores con cardinalidad masiva arbitraria (`nameOrig`, `nameDest`).
   * Extracción de semántica del destinatario: nueva variable `dest_type` (`C` = Cliente particular, `M` = Comercio).
   * Eliminación de la bandera estática `isFlaggedFraud` (que sólo detectaba el 0.19% de los fraudes reales).
3. **Tratamiento de Valores Atípicos (*Outliers*):**
   * Recorte (*winsorizing*) sobre `amount` en el percentil 99.99 ($9,615,000 USD), mitigando inestabilidad numérica en gradientes.
4. **Ingeniería de Características (*Feature Engineering*):**
   * **Discrepancias contables:** Cálculo de balances esperados vs reales:
     $$\text{errorBalanceOrig} = \text{newbalanceOrig} + \text{amount} - \text{oldbalanceOrg}$$
     $$\text{errorBalanceDest} = \text{oldbalanceDest} + \text{amount} - \text{newbalanceDest}$$
   * **Variables temporales periódicas:** A partir de `step` (horas acumuladas):
     * `hora_del_dia` ($step \pmod{24}$)
     * `dia_de_la_semana` ($(step // 24) \pmod 7$)
5. **Codificación Categórica:**
   * One-Hot Encoding binario (`pd.get_dummies` con `drop_first=True`) para `type` y `dest_type`, preservando tipo entero en memoria.
6. **Estrategia de Balance de Clases:**
   * Ponderación analítica de la función de pérdida (`compute_class_weight` balanceado) asignando pesos inversamente proporcionales a las frecuencias (Clase 0 $\approx 0.5006$, Clase 1 $\approx 387.35$), evitando sobremuestreo sintético (SMOTE) para minimizar la sobrecarga en las etapas concurrentes de Go.
7. **División Cronológica y Estandarización:**
   * Ordenamiento temporal estricto por `step`.
   * Split 80/20: **5,090,096 transacciones de entrenamiento** y **1,272,524 de prueba**, previniendo fuga de datos temporal (*temporal data leakage*).
   * Estandarización `StandardScaler` ajustada sobre entrenamiento y aplicada a prueba sobre las **14 variables predictoras finales**.

---

## 📁 4. Estructura del Repositorio

```text
├── limpieza_paysim_fraud.ipynb     # Notebook con EDA, limpieza y exportación
├── paysim.csv                      # Dataset masivo (>6.36M filas, Git LFS)
├── data-limpia/                    # Muestras y salidas del preprocesamiento
│   ├── paysim_sample_50k.csv       # Muestra de 50,000 registros procesados
│   ├── y_train_sample.csv          # Vector de etiquetas de entrenamiento
│   └── y_test_sample.csv           # Vector de etiquetas de prueba
├── .gitattributes                  # Filtros de Git LFS para paysim.csv
├── .gitignore                      # Configuración de exclusión de Git
└── README.md                       # Documentación del proyecto
```
