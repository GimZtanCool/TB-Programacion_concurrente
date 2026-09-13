# Detección de Fraude en Transacciones Financieras (IEEE-CIS Fraud Detection)
### CC65 - Programación Concurrente y Distribuida | Trabajo Parcial 2026-20

Repositorio del proyecto del curso **Programación Concurrente y Distribuida (CC65)**, enfocado en el procesamiento masivo de datos, análisis exploratorio (EDA) y preparación de un conjunto de datos a gran escala para el entrenamiento e inferencia de modelos de Machine Learning mediante arquitecturas concurrentes y paralelas (Go, Goroutines, Worker Pools y verificación formal en Promela/Spin).

---

## 📌 1. Contexto y Objetivos del Proyecto

El fraude en transacciones electrónicas representa pérdidas multimillonarias anuales a nivel global y vulnera la confianza en la infraestructura financiera digital. Los objetivos del proyecto son:
1. **Analizar y preprocesar** un conjunto de datos masivo (>590,000 transacciones con más de 400 atributos) garantizando consistencia, balanceo y escalabilidad.
2. **Implementar modelos de Machine Learning** bajo esquemas secuenciales y concurrentes en lenguaje **Go** (sin librerías de terceros para la concurrencia), explotando mecanismos nativos como `goroutines`, canales (`chan`), `sync.WaitGroup` y `sync.Mutex`.
3. **Evaluar el rendimiento computacional** mediante métricas de *Speedup*, escalabilidad y medias recortadas, modelando la sincronización en **Promela/Spin** para verificar formalmente la ausencia de condiciones de carrera y *deadlocks*.

### Alineación con los ODS (Objetivos de Desarrollo Sostenible)
* **ODS 8: Trabajo decente y crecimiento económico:** Fortalecimiento de la seguridad en sistemas de pago digitales e inclusión financiera confiable.
* **ODS 9: Industria, innovación e infraestructura:** Construcción de infraestructura de software robusta, concurrente y de alta disponibilidad para el procesamiento de transacciones en tiempo real.
* **ODS 16: Paz, justicia e instituciones sólidas:** Reducción sustancial del flujo financiero ilícito y actividades fraudulentas en el comercio electrónico.

---

## 📊 2. Dataset Seleccionado: IEEE-CIS Fraud Detection

El dataset proviene de la competencia oficial organizada por el **IEEE Computational Intelligence Society (IEEE-CIS)** y **Vesta Corporation**:
* **Fuente:** [IEEE-CIS Fraud Detection en Kaggle](https://www.kaggle.com/c/ieee-fraud-detection)
* **Estructura original:**
  * `train_transaction.csv`: 590,540 registros y 394 columnas (montos, tarjetas, temporizadores, atributos anonimizados).
  * `train_identity.csv`: 144,233 registros y 41 columnas (datos de red, navegador, identidad del dispositivo).
* **Unión (*Left Join*):** Sobre la clave `TransactionID`, consolidando una matriz inicial de **590,540 transacciones y 434 columnas**.
* **Variable objetivo:** `isFraud` (binaria: 0 = legítima, 1 = fraude).
  * **Distribución:** 569,877 transacciones legítimas (~96.5%) vs. 20,663 fraudulentas (~3.5%), reflejando el severo desbalance de clases inherente al problema.

---

## 🛠️ 3. Pipeline de Limpieza y Preprocesamiento (`limpieza_ieee_cis_fraud.ipynb`)

El procedimiento implementado en el notebook sigue un flujo estricto para evitar fuga de información (*data leakage*):

1. **Tratamiento de Valores Faltantes:**
   * Diagnóstico: 49.3% de las columnas presentaba más del 50% de valores nulos.
   * Se descartaron **12 columnas** con más del **90% de nulos** (principalmente `id_*` y `dist2`).
   * Imputación de variables numéricas restantes con la **mediana**.
   * Imputación de variables categóricas con el valor explícito `"missing"`.
2. **Tratamiento de Valores Atípicos (*Outliers*):**
   * Recorte (*winsorizing*) sobre `TransactionAmt` en el percentil 99.99 ($4,222.95 USD), eliminando 57 registros no representativos.
3. **Codificación de Variables Categóricas:**
   * **Alta cardinalidad** (`DeviceInfo`, `P_emaildomain`, `R_emaildomain`, etc.): agrupación de frecuencias menores al 1% en la etiqueta `"rare"`.
   * **Baja cardinalidad**: transformación con **One-Hot Encoding** (`pd.get_dummies` con `drop_first=True`), preservando tipo entero (`dtype=int`).
4. **Ingeniería de Características Temporales:**
   * A partir de `TransactionDT` (segundos transcurridos), se derivaron variables periódicas de comportamiento: `hora_del_dia` (0-23) y `dia_de_la_semana` (0-6).
5. **Reducción de Dimensionalidad (PCA):**
   * El bloque de atributos propietarios de Vesta (`V1` a `V339`) presentaba alta redundancia y multicolinealidad.
   * Se aplicó **PCA** conservando el 95% de la varianza explicada, colapsando 339 columnas continuas en únicamente **2 componentes principales** (`V_pca_1`, `V_pca_2`).
6. **Balance de Clases:**
   * Cálculo de pesos balanceados (`class_weight`: Clase 0 $\approx 0.518$, Clase 1 $\approx 14.289$) para ponderar la función de pérdida sin recurrir a técnicas de sobremuestreo sintético (SMOTE), reduciendo el overhead en las etapas concurrentes de Go.
7. **División Temporal y Estandarización:**
   * Ordenamiento cronológico por `TransactionDT`.
   * Split temporal: **80% entrenamiento** (472,386 registros) y **20% prueba** (118,097 registros).
   * Estandarización (`StandardScaler`) ajustada exclusivamente sobre el conjunto de entrenamiento y aplicada sobre el conjunto de prueba para las **145 variables predictoras finales**.

---

## 📁 4. Estructura del Repositorio

```text
├── CC65_PCs_TP-202620.pdf          # Pautas y rúbrica oficial del curso
├── limpieza_ieee_cis_fraud.ipynb   # Notebook con EDA, limpieza y exportación
├── data-limpia/                    # Salida del preprocesamiento
│   ├── y_train.csv                 # Etiquetas de entrenamiento (472,386 filas)
│   ├── y_test.csv                  # Etiquetas de prueba (118,097 filas)
│   ├── X_train.csv                 # Matriz de entrenamiento escalada (generada localmente)
│   └── X_test.csv                  # Matriz de prueba escalada (generada localmente)
├── .gitattributes                  # Filtros de Git LFS para archivos de datos
├── sample_submission.csv           # Plantilla de predicción de la competencia
├── train_identity.csv              # Datos de identidad originales (Git LFS)
├── train_transaction.csv           # Datos de transacciones originales (Git LFS)
└── README.md                       # Documentación del proyecto
