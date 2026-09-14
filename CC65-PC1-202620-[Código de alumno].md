![](data:image/jpeg;base64...)

UNIVERSIDAD PERUANA DE CIENCIAS APLICADAS

FACULTAD DE INGENIERÍA

CIENCIAS DE LA COMPUTACIÓN

CURSO: 1ACC0065 – Programación Concurrente y Distribuida

DOCENTE: Herminio Paucar Curasma

NRC: 8646

TB1

INTEGRANTES:

Avalos Sánchez, César Gabriel (u202310307) – Ciencias de la Computación

Rojas Cuadros, Fabian Marcelo (u202218498) – Ciencias de la Computación

Ballon Villar, Diego Eduardo (u201520327) – Ciencias de la Computación

Cuadrado Jimenez, Sebastian Antonio (u202312882) – Ciencias de la Computación

Ciclo Académico: 202620

Indice

[1. Resumen del trabajo 3](#_jk84ejt0rcco)

[2. Objetivos del trabajo 3](#_iqerd3wh8pv4)

[3. Explicación de cada tema investigado 4](#_fjcee7lz1qks)

[4. Análisis del caso de uso. 10](#_xytq88yqzom7)

[5. Explicación del procedimiento de limpieza del dataset. 12](#_srp7qbdpwmaw)

[● Descripción general del conjunto de datos 12](#_b71nxswt0v40)

[● Integridad y diagnóstico de valores faltantes y duplicados 13](#_diag_nulos)

[● Tratamiento de variables identificadoras y regla heurística 13](#_trat_ids)

[● Tratamiento de valores atípicos en montos 14](#_trat_outliers)

[● Ingeniería de características contables y temporales 14](#_feat_eng)

[● Codificación de variables categóricas 15](#_cod_cat)

[● Balance de clases mediante ponderación 15](#_bal_clases)

[● División del conjunto de datos (Split temporal) 16](#_div_datos)

[● Estandarización 16](#_estandarizacion)

[Tabla 1. Resumen del procedimiento de limpieza y preprocesamiento 16](#_tabla1)

[6. Trabajo en github, historial de commits de cada integrante del grupo. 17](#_k6jfl926m5g7)

[7. Referencias. 18](#_lhqbiid7mfob)

# **Resumen del trabajo**

El presente trabajo forma parte de la Práctica Calificada 1 (PC1) del curso de Programación Concurrente y Distribuida, y tiene como objetivo investigar cómo la programación concurrente y paralela potencia el rendimiento de modelos de machine learning al enfrentar grandes volúmenes de datos. Para ello, el grupo seleccionó y analizó cuatro papers académicos publicados en los últimos años (2023-2025), obtenidos de fuentes indexadas como Cambridge Data & Policy, IEEE Xplore, MDPI y arXiv, accedidos a través del Centro de Información de la UPC.

Los papers investigados abordan distintos casos de uso de clasificación mediante machine learning: (1) la priorización de supervisión en contratación pública mediante modelos de regresión logística y Random Forest (VigIA), (2) la detección de fraude financiero en transacciones digitales utilizando algoritmos de gradient boosting como XGBoost y LightGBM, (3) la aceleración del entrenamiento de modelos de ensamble (Random Forest, AdaBoost, XGBoost, KNN) mediante procesamiento paralelo en entornos multicore, y (4) el análisis del consumo energético asociado a distintos paradigmas de paralelización (memoria compartida, memoria distribuida y GPU) en el entrenamiento de modelos de ML.

A partir de este análisis bibliográfico, el grupo seleccionó como caso de uso propio la detección de fraude en transacciones financieras, tomando como fuente el dataset PaySim1 (Synthetic Financial Datasets For Fraud Detection, Kaggle), en línea con el paper investigado por César Gabriel Avalos Sánchez ("Financial Fraud Detection using Machine Learning Models"), que aborda directamente esta problemática utilizando el mismo dataset. Sobre este dataset, compuesto por más de 6.3 millones de registros de transacciones simuladas, se aplicó un procedimiento de limpieza y preprocesamiento que incluyó tratamiento de variables identificadoras, tratamiento de valores atípicos, codificación de variables categóricas, ingeniería de variables sobre los saldos de las cuentas y balanceo de clases. Todo el desarrollo del trabajo, incluyendo el historial de commits de cada integrante, fue gestionado a través de un repositorio de GitHub siguiendo la metodología Git Flow.

# **Objetivos del trabajo**

**Objetivo general:** Investigar y analizar el uso de la programación concurrente y distribuida en modelos de machine learning, evaluando los beneficios que esta aporta en términos de eficiencia computacional, tiempos de entrenamiento e inferencia, mediante la revisión de literatura académica reciente y la aplicación práctica sobre un dataset de gran volumen (>1,000,000 de registros).

**Objetivos específicos:**

* Identificar y analizar tres papers académicos por integrante del grupo (2021-2025), publicados en fuentes indexadas reconocidas (Scopus, IEEE Xplore, ACM Digital Library, entre otras), relacionados con el uso de concurrencia en modelos de machine learning.
* Describir los algoritmos, frameworks, bibliotecas y patrones de concurrencia empleados en cada paper investigado, incluyendo su pseudocódigo cuando corresponda.
* Emitir una opinión crítica y fundamentada sobre los beneficios y/o perjuicios que aporta la programación concurrente en cada uno de los casos investigados.
* Seleccionar y sustentar un caso de uso propio, justificando la elección de un dataset con más de un millón de registros.
* Aplicar un procedimiento de limpieza y preprocesamiento de datos sobre el dataset seleccionado, documentando cada paso realizado.
* Gestionar el desarrollo colaborativo del trabajo mediante un repositorio de GitHub, evidenciando el historial de commits individual de cada integrante del grupo.

# **Explicación de cada tema investigado**

| **Alumno** | **Diego Eduardo Ballon Villar** | **César Gabriel Avalos Sánchez** | **Fabian Marcelo Rojas Cuadros** | **Sebastian Cuadrado** |
| --- | --- | --- | --- | --- |
| **Título del paper** | VigIA: Prioritizing Public Procurement Oversight with Machine Learning Models and Risk Indices | Financial Fraud Detection using Machine Learning Models | A Parallel Approach to Enhance the Performance of Supervised Machine Learning Realized in a Multicore Environment | Performance and Energy Consumption of Parallel Machine Learning Algorithms |
| **Autores** | Andrés Salazar, Juan F. Pérez, Jorge Gallego | Anita Ukwu, Shatha Ghareeb, Jamila Mustafina Kazab | Ashutosh Ghimire y Fathi Amsaad - Department of Computer Science and Engineering, Wright State University, Dayton, OH, EE.UU. | Xidong Wu, Preston Brazzle, Stephen Cahoon |
| **Año** | 2024 | 2025 | 2024 | 2023 |
| **Motivación** | Los organismos de control de contratación pública tienen presupuesto y personal limitados frente a un volumen de contratos que hace inviable revisar cada proceso individualmente. El caso real: la Veeduría Distrital de Bogotá, con un presupuesto anual equivalente a apenas el 0,1% del presupuesto total de la ciudad, debe vigilar cientos de miles de contratos cada año. | El crecimiento constante de los sistemas de pago digitales y transacciones con tarjeta de crédito genera un volumen masivo de datos expuestos a ataques de fraude exigiendo modelos capaces de procesar información rápidamente. | El crecimiento en tamaño y complejidad de los datasets hace que entrenar modelos de ML consuma cada vez más tiempo. Los autores buscan acelerar el entrenamiento aprovechando CPUs multicore (hardware accesible, sin necesidad de GPU), aplicándolo a un problema de clasificación real: predecir la calidad de un vino a partir de sus propiedades fisicoquímicas. | El entrenamiento de modelos de *Machine Learning* exige grandes conjuntos de datos y múltiples iteraciones computacionales. Sin embargo, la gran mayoría de las investigaciones sobre el entrenamiento concurrente se enfocan únicamente en el rendimiento, dejando de lado el consumo de energía, que es una métrica fundamental para aplicaciones de alto rendimiento. |
| **Problema que resuelve** | Cómo priorizar, entre un volumen alto de contratos, cuáles supervisar primero, cubriendo dos tipos de desperdicio a la vez: el activo (corrupción deliberada) y el pasivo (ineficiencia por falta de capacidad), que según la literatura previa representan hasta el 83% del desperdicio total en algunas entidades estatales. También aborda el "problema de etiquetado selectivo": los modelos entrenados solo con casos ya detectados por la Contraloría heredan el sesgo de lo que esa entidad ya encontró. | Este trabajo resuelve la necesidad de comparar, optimizar y seleccionar la arquitectura de machine learning que ofrezca el mejor equilibrio entre precision, recall y eficiencia computacional, permitiendo reducir drasticamente los falsos positivos e identificar patrones complejos de ciberataques sin comprometer los tiempos de respuesta exigidos en la industria bancaria. | Clasificar la calidad de un vino (escala de 3 a 9) a partir de 11 atributos fisicoquímicos (acidez, alcohol, pH, azúcar residual, sulfatos, etc.) es una tarea de clasificación multiclase con un dataset pequeño (4,898 instancias) pero con 7 clases de salida, lo que la hace más compleja que un problema binario. El reto es lograr que el entrenamiento paralelo siga siendo eficiente incluso cuando el dataset es chico, donde el overhead de paralelizar puede no compensar la ganancia. | La falta de análisis y optimización sobre cómo impacta la paralelización no solo en el tiempo de ejecución computacional, sino en la eficiencia energética cuando se aplican algoritmos de *machine learning* en plataformas de alto rendimiento. |
| **Propuesta de solución** | Dos componentes complementarios, ambos a nivel de contrato: (1) modelos de machine learning que predicen sobrecostos y demoras usando variables objetivas del propio proceso; (2) un índice compuesto de riesgo de irregularidad (IRIC), construido como el promedio de 11 variables binarias en tres dimensiones: falta de competencia, falta de transparencia y anomalías del proceso. Incluye una versión ponderada por el valor del contrato (IRICP). | Evaluación comparativa de algoritmos de clasificación (Logistic Regression, Naive Bayes, Random Forest, LightGBM y XGBoost), aprovechando la escalabilidad de los modelos basados en árboles impulsados por gradiente (gradient boosting) | Aplicar la misma arquitectura paralela de ensamble (Random Forest, XGBoost, AdaBoost, KNN) al dataset de calidad de vino (Wine Quality Dataset, UCI), particionando el entrenamiento entre los cores disponibles, y comparando el resultado contra el mismo experimento con el dataset de fraude (mucho más grande) para observar cómo cambia el beneficio de la concurrencia según el tamaño del dataset. | Implementar varios algoritmos clásicos de aprendizaje automático bajo diferentes paradigmas de paralelización (memoria compartida, memoria distribuida y aceleración por GPU) para medir empíricamente cómo la complejidad del modelo y el tamaño de los datos de entrenamiento afectan la eficiencia paralela y la energía consumida. |
| **Algoritmos usados y Pseudocodigo** | Algoritmo A: Cálculo del Índice de Riesgo de Irregularidad en Contratación (IRIC) Entrada: Registro de contrato c, historial del proveedor Salida: IRIC(c) e IRICP(c)  1. Para cada contrato c:  2. Calcular 11 variables binarias (0/1) en 3 dimensiones:  3. Falta de competencia (5 variables): proponente único, proveedor multipropósito, historial excepcionalmente alto del proveedor, contratación directa, régimen especial, período de aviso extremo  4. Falta de transparencia (2 variables): errores o datos faltantes, período de decisión extremo  5. Anomalías del proceso (3 variables): proveedor con sobrecostos previos, proveedor con retrasos previos, ausencia de proceso de contratación  6. IRIC(c) = promedio aritmético de las 11 variables binarias  7. peso(c) = log(valor\_contrato(c)) / percentil\_75(log(valor), tipo\_contrato)  8. IRICP(c) = IRIC(c) × peso(c)  9. Retornar IRIC(c), IRICP(c) Algoritmo B: Modelo de predicción de ineficiencias (sobrecostos/retrasos) Entrada: Dataset de contratos con variables del proceso Salida: Modelo entrenado + variables más relevantes  1. Balancear clases mediante remuestreo (igual número de contratos con y sin sobrecostos/retrasos)  2. Estandarizar variables explicativas  3. Entrenar Random Forest sobre todas las variables disponibles  4. Calcular importancia de características (feature importance) por árbol  5. Seleccionar subconjunto reducido de variables más importantes (ej. 7 variables para sobrecostos, 6 para retrasos)  6. Entrenar modelo final con Regresión Logística y Random Forest usando el subconjunto reducido (70% train / 30% test) 7. Calcular gráficos de dependencia parcial por variable 8. Evaluar con exactitud, precisión, recall, curva ROC y AUC 9. Retornar modelo con mejor AUC (Random Forest, AUC=0.943 vs 0.821 de Regresión Logística) | Algoritmo: Pipeline de detección de fraude con balanceo SMOTE Entrada: Dataset de transacciones (100,000 registros, 1% fraude) Salida: Modelo entrenado y evaluado  1. Cargar dataset y verificar ausencia de nulos/duplicados 2. Codificar variables categóricas (TransactionType, Location, MerchantID) mediante label encoding 3. Normalizar variables numéricas (transaction amount) con MinMaxScaler  4. Aplicar SMOTE (Synthetic Minority Over-Sampling Technique):  5. Para cada instancia de la clase minoritaria (fraude):  6. Generar muestras sintéticas interpolando entre vecinos cercanos de la misma clase  7. Balancear el dataset con las muestras sintéticas generadas  8. Dividir dataset en 80% entrenamiento / 20% prueba (muestreo estratificado para mantener proporción fraude/no-fraude)  9. Para cada modelo en [Regresión Logística, Naive Bayes, Random Forest, XGBoost, LightGBM]:  10. Aplicar GridSearchCV con validación cruzada estratificada de 5 folds para optimizar hiperparámetros 11. Entrenar modelo con datos balanceados  12. Evaluar con exactitud, precisión, recall, F1-score y ROC-AUC  13. Comparar los 5 modelos (+ Autoencoder como baseline de detección de anomalías)  14. Seleccionar mejor modelo (Random Forest: 99% exactitud, 0.99 recall, AUC=1.00) | Algoritmo 1: Bootstrap Sampling para generar subconjuntos Entrada: Dataset original T (tamaño N), número de subconjuntos M Salida: Lista de subconjuntos S  1. Inicializar lista vacía S  2. Para i = 1 hasta M:  3. Inicializar subconjunto vacío subset\_i  4. Mientras tamaño(subset\_i) < N:  5. Seleccionar aleatoriamente una instancia de T (con reemplazo)  6. Agregar la instancia a subset\_i  7. Agregar subset\_i a S  8. Retornar S Algoritmo: Entrenamiento paralelo del ensamble (Random Forest)  1. Particionar datos de entrenamiento con Bootstrap Sampling (Algoritmo 1)  2. Para cada subconjunto (en paralelo, uno por core):  3. Construir un árbol de decisión (CART) sobre ese subconjunto  4. Seleccionar subconjunto aleatorio de features en cada split  5. Combinar todos los árboles entrenados (Ensemble)  6. Predicción final = voto mayoritario de todos los árboles | Algorithm 1 (del paper): Stochastic Gradient Descent (SGD) Entrada: Número de épocas de entrenamiento T, tasa de aprendizaje η, tamaño de mini-batch b Inicializar: Parámetros del modelo x0  1. Para t = 1, 2, ..., T:  2. Extraer muestras de mini-batch Bt = {ξ^j}, j=1..b, con |Bt| = b, del dataset de entrenamiento D  3. xt = xt-1 - η∇xf(xt-1; Bt)  4. Retornar xT Paralelización síncrona (shared-memory, con OpenMP):  5. Usar "#pragma omp parallel for" para acelerar la multiplicación de vectores en cada iteración  6. No paralelizar por batch-size, para evitar condiciones de carrera en la actualización del modelo Paralelización asíncrona (shared-memory):  7. Usar "#pragma omp parallel for schedule(static,1)" para paralelizar cada iteración de actualización  8. Usar "#pragma omp critical" para forzar que los hilos actualicen el modelo uno a la vez (evitar condiciones de carrera) Paralelización distribuida (MPI, síncrona):  9. Cada nodo cliente muestrea una porción de los datos y calcula su propio gradiente local  10. Un servidor central agrega los gradientes de todos los clientes  11. El servidor actualiza el modelo y lo redistribuye a todos los nodos Paralelización distribuida (MPI, asíncrona):  12. Al recibir un gradiente de un nodo, el servidor actualiza el modelo inmediatamente y lo devuelve solo a ese nodo (sin sincronización global) |
| **Descripción de servicios, framework, bibliotecas usadas, código, lenguaje de programación, patrones.** | Regresión logística y Random Forest como modelos de clasificación (elegidos explícitamente por su interpretabilidad frente a redes neuronales, que se descartan por su opacidad). Selección de variables mediante la importancia de las características del Random Forest. Gráficos de dependencia parcial para explicar cada predictor. Métricas de evaluación: exactitud, precisión, exhaustividad, curva ROC y AUC. Fuente de datos: la plataforma SECOP II de Colombia Compra Eficiente, cruzada con el registro de sanciones de la Superintendencia de Industria y Comercio, ambas de datos.gov.co. No especifican el lenguaje de programación en el cuerpo del artículo, pero el stack de modelos (regresión logística, random forest, curvas ROC) es estándar en Python con scikit-learn o en R con caret/randomForest. | Uso de pandas y Numpy para la carga, filtrado y transformación matricial y preprocesamiento de vectores de características en memoria.  Implementación de la biblioteca scikit-learn para la construcción de los modelo base (logistic regression, naive bayes, random forest), la codificación de variables categóricas, el escalamiento de datos y la división de conjuntos de entrenamiento y prueba.  Utilización de XGboost y LightGBM para la paralelización multihilo nativa a nivel del cpu durante la búsqueda de divisiones óptimas de los árboles y el cálculo de matrices de hessiana y gradiente.  Aplicación de métricas para conjuntos desbalanceados mediante scikit-learn.metrics, incluyendo precisión, recall, f1 y análisis de curvas ROC y AUC. Synthetic Financial Datasets For Fraud Detection (PaySim1) <https://www.kaggle.com/datasets/ealaxi/paysim1> | Python 3 (en los 3 dispositivos de prueba). **Preprocesamiento:** PCA (Análisis de Componentes Principales) para reducir dimensionalidad del dataset de fraude (ya venía aplicado en el dataset original de Kaggle) y estandarización (media 0, desviación 1) en ambos datasets; split 80/20 train/test. **Patrón de diseño:** paralelismo de datos vía data partitioning + bootstrap sampling, con sincronización entre procesos/cores para agregar resultados (patrón map-reduce a nivel conceptual: particionar -> procesar en paralelo -> combinar/votar). **Código fuente disponible en GitHub:**<https://github.com/aashutoshghimire/parallel-ensemble-model>. **Datasets:** Kaggle (Fraud Detection Credit Card) y UCI ML Repository (Wine Quality). | Implementaciones desarrolladas en los lenguajes de programación C++ (para regresión logística y algoritmos genéticos) y Python (para redes neuronales). Utiliza patrones de paralelismo de memoria compartida, paralelismo de memoria distribuida y delegación de procesamiento mediante aceleración por GPU. |
| **Opinión crítica de su parte (estudiante) de los beneficios o perjuicios del uso de la programación concurrente.** | **Fabian:** El paper no menciona concurrencia de forma explícita, pero el uso de Random Forest ya trae ese potencial implícito: cada árbol del ensamble puede entrenarse de forma independiente, algo que confirma el paper de Multicore que yo investigué. El perjuicio es que, al priorizar interpretabilidad sobre rendimiento (descartan redes neuronales por opacidad), los autores no exploraron si ese mismo Random Forest podría escalar a nivel nacional sin rediseñar el pipeline para procesamiento paralelo.  **César:** Comparado con mi paper de Financial Fraud, donde XGBoost y LightGBM sí discuten su paralelismo interno, VigIA se queda corto en ese aspecto: usan Random Forest pero nunca hablan de cuántos núcleos o qué tiempo de entrenamiento les tomó procesar los datos de SECOP II. Dado que manejan "cientos de miles de contratos", el beneficio de paralelizar sería real, pero el paper simplemente no lo cuantifica, a diferencia de otros papers del grupo que sí miden el speedup.  **Sebastian:** Desde la óptica de mi paper (que mide consumo energético del paralelismo), noto que VigIA tampoco discute el costo computacional de calcular los gráficos de dependencia parcial para cada predictor, algo que sí podría beneficiarse de cómputo paralelo si se aplica a los cientos de miles de contratos mencionados. Es un vacío compartido con otros papers del grupo: se preocupan por la exactitud del modelo, pero ignoran completamente el costo energético o de tiempo de escalarlo. | **Diego:** Este paper sí es más explícito que el mío en el uso de concurrencia: XGBoost y LightGBM paralelizan de forma nativa el cálculo de gradientes y la búsqueda de splits en los árboles. El beneficio es claro para datos financieros masivos en tiempo real. Sin embargo, no explican cómo manejan la sincronización entre hilos cuando los datos están tan desbalanceados (fraude vs. no fraude), lo cual podría generar particiones de trabajo muy desiguales entre cores.  **Fabian:** Comparado con el paper de Multicore que yo trabajé, que sí mide el speedup exacto por número de cores, este paper de fraude se queda en un nivel más superficial: menciona que XGBoost y LightGBM paralelizan internamente, pero no reporta tiempos de entrenamiento ni cuántos hilos usaron. Es un beneficio "de caja negra", confían en que la librería ya resuelve la concurrencia, sin analizar si esa paralelización interna es óptima para su caso de uso específico.  **Sebastian:** Desde la perspectiva energética que trabajé en mi paper, este es el caso donde más se nota la ausencia de esa métrica: paralelizar árboles de gradiente sobre "volúmenes masivos de datos financieros en tiempo real" tiene un costo energético que ni siquiera se menciona. El perjuicio que señalan (gestión cuidadosa de memoria para evitar cuellos de botella) va en la línea correcta, pero se queda corto al no conectar ese cuello de botella con el consumo energético que representa. | **Diego:** Este es el paper más riguroso del grupo en términos de medir concurrencia: cuantifica el speedup real por número de cores (hasta 9.05× para Random Forest). Comparado con VigIA, que ni siquiera discute el tema, aquí queda clarísimo que no todos los algoritmos se benefician igual del paralelismo, KNN apenas mejora, mientras que los basados en árboles casi escalan linealmente. El perjuicio es que todo el experimento se queda en un solo nodo multicore, sin explorar si las conclusiones se sostendrían en un clúster distribuido como el que necesitaría VigIA a escala nacional.  **César:** A diferencia de mi paper de fraude financiero, que asume la paralelización de XGBoost como una caja negra, este paper sí abre esa caja: explica el bootstrap sampling, la construcción paralela de árboles y la combinación por votación. El beneficio es una comprensión mucho más profunda del "cómo". La limitación que veo es que el dataset de fraude que usan (284,807 registros) es mucho menor al de mi paper, así que sus conclusiones sobre "datasets grandes" podrían no sostenerse con volúmenes de datos verdaderamente masivos (millones de registros).  **Sebastian:** Este paper comparte con el mío el interés en medir el paralelismo de forma empírica, pero se queda solo en tiempo de ejecución, sin tocar el consumo energético que sí exploro yo. Sería interesante cruzar ambos hallazgos: si Random Forest escala casi linealmente en velocidad con más cores, ¿escala también linealmente en consumo de energía, o hay ahí también rendimientos decrecientes? El paper no da pistas sobre eso, y es una oportunidad perdida para complementar su análisis de "punto de quiebre" del paralelismo | **Diego:** Este paper aporta algo que ningún otro paper del grupo considera: el consumo energético, no solo la velocidad. Frente a VigIA, que ni siquiera piensa en escalabilidad, este trabajo obliga a preguntarse si vale la pena paralelizar un modelo si el ahorro de tiempo no compensa el gasto energético adicional, una pregunta que sería muy relevante aplicar al caso de contratación pública si algún día se lleva a producción a gran escala.  **César:** Comparado con mi paper de fraude, que usa paralelismo interno de librerías (XGBoost/LightGBM) sin cuestionarlo, este paper analiza explícitamente distintos paradigmas (memoria compartida, memoria distribuida, GPU), lo cual es mucho más completo. El beneficio es que da un panorama más realista de las opciones de concurrencia disponibles. El perjuicio es que, al combinar C++ para regresión logística/algoritmos genéticos y Python para redes neuronales, la comparación de consumo energético entre lenguajes podría estar sesgada por las diferencias inherentes de eficiencia entre ambos, no solo por el paradigma de paralelismo usado.  **Fabian:** Comparado con el paper de Multicore que yo investigué, que mide speedup pero ignora energía, este paper llena ese vacío de forma directa. El beneficio es evidente: en aplicaciones de alto rendimiento, el consumo energético es tan crítico como el tiempo de cómputo, algo que mi propio paper no consideró. El perjuicio que noto es que no comparan sus resultados de energía contra los mismos algoritmos (Random Forest, XGBoost) que usa el paper de Multicore, lo cual hubiera permitido una comparación directa entre "cuánto se gana en tiempo" y "cuánto se gasta en energía" para los mismos modelos. |
| **Link de acceso al paper.** | <https://www.cambridge.org/core/journals/data-and-policy/article/vigia-prioritizing-public-procurement-oversight-with-machine-learning-models-and-risk-indices/34D04747A94A7099E3CD8B91221338ED> | <https://ieeexplore.upc.elogim.com/document/11367919> | <https://www.mdpi.com/2504-4990/6/3/90> | <https://arxiv.org/abs/2305.00798> |

# **Análisis del caso de uso.**

* Contexto y problema abordado

El estudio parte de un problema del sector financiero: el crecimiento sostenido de los sistemas de pago digitales y las transacciones con tarjeta genera un volumen masivo de datos expuestos a fraude, lo que exige modelos capaces de procesar información rápidamente y con alta precisión. El desafío no es solo detectar fraude, sino hacerlo minimizando los falsos positivos —que generan fricción innecesaria al usuario— sin sacrificar los tiempos de respuesta que exige la industria bancaria..

* Justificación de la elección del dataset: Synthetic Financial Datasets For Fraud Detection (PaySim1)

La elección de PaySim1 como dataset propio del equipo responde a tres criterios. Primero, cumple con holgura el requisito de volumen del curso: sus 6.3 millones de registros son varias veces superiores al mínimo de un millón exigido, lo que garantiza que las técnicas de procesamiento concurrente y distribuido evaluadas tengan un impacto medible en tiempos de entrenamiento. Segundo, mantiene coherencia directa con la revisión bibliográfica del grupo, ya que reproduce el mismo dataset utilizado en el paper de Ukwu et al. (2025) analizado por César Avalos, lo que permite contrastar los resultados propios con los reportados en la literatura y evaluar si las afirmaciones del paper sobre paralelización nativa en XGBoost/LightGBM se sostienen empíricamente. Tercero, su fuerte desbalance de clases (apenas 0.13% de transacciones fraudulentas) y su estructura transaccional con componente temporal (step) lo convierten en un caso de uso representativo de los desafíos reales de detección de fraude en sistemas financieros a gran escala, exigiendo no solo modelos precisos sino también eficientes en su entrenamiento e inferencia.

● Objetivo y enfoque técnico

El objetivo es comparar, optimizar y seleccionar la arquitectura de machine learning que ofrezca el mejor equilibrio entre precisión, recall y eficiencia computacional. Para esto se evalúan cinco algoritmos de clasificación —Regresión Logística, Naive Bayes, Random Forest, XGBoost y LightGBM— sobre el dataset sintético PaySim1 (Synthetic Financial Dataset For Fraud Detection).

El pipeline sigue un flujo estándar: carga y preprocesamiento con pandas/NumPy, codificación de variables categóricas y escalamiento, división 80/20 en entrenamiento y prueba, entrenamiento de cada modelo y evaluación final con precisión, recall, F1 y curva ROC/AUC.

El componente de concurrencia aparece de forma específica en XGBoost y LightGBM, que aprovechan multihilo nativo a nivel de CPU durante dos operaciones puntuales del boosting: el cálculo en paralelo del gradiente y la hessiana de la función de pérdida en cada nodo del árbol, y la búsqueda paralela del mejor split entre hilos. A diferencia de un esquema de bagging (como Random Forest), donde los árboles se construyen de forma independiente, aquí los árboles se combinan secuencialmente (boosting), por lo que la paralelización ocurre dentro de la construcción de cada árbol, no entre árboles.

● Resultados relevantes

El paper reporta métricas de calidad del modelo (precisión, recall, F1, ROC-AUC) para comparar los cinco algoritmos y así "seleccionar el modelo con mejor equilibrio precision/recall/eficiencia". Sin embargo, no se reportan tiempos de entrenamiento, número de hilos usados, ni ningún tipo de medición de speedup asociado a la paralelización de XGBoost/LightGBM. La ganancia de eficiencia computacional se afirma pero no se cuantifica empíricamente.

● Aporte del trabajo

El valor del paper está en la comparación sistemática de algoritmos clásicos y de gradient boosting para un caso de uso financiero real (PaySim1), evaluados bajo métricas adecuadas para datos desbalanceados (recall, F1, AUC en lugar de solo accuracy). El uso de XGBoost y LightGBM se justifica precisamente por su escalabilidad y su paralelización nativa a nivel de CPU, lo cual los hace atractivos para "volúmenes masivos de datos financieros en tiempo real".

● Limitaciones

Paralelismo de "caja negra": el paper menciona que XGBoost y LightGBM paralelizan internamente el cálculo de gradiente/hessiana y la búsqueda de splits, pero no analiza si esa paralelización es óptima para su caso de uso específico, ni reporta cuántos hilos se usaron ni el tiempo de entrenamiento resultante. Se confía en que la librería resuelve la concurrencia sin verificarlo empíricamente.

Sincronización bajo desbalance de clases: no se explica cómo se maneja la sincronización entre hilos cuando los datos están fuertemente desbalanceados (fraude vs. no fraude), lo que podría generar particiones de trabajo desiguales entre cores durante la búsqueda de splits.

Ausencia de métrica energética: el paper menciona la necesidad de una "gestión cuidadosa de memoria para evitar cuellos de botella", pero no conecta ese cuello de botella con el consumo energético que implica paralelizar árboles de gradiente sobre datos masivos en tiempo real.

Sin comparación de tiempos entre modelos: al no reportar tiempos de entrenamiento para ningún algoritmo (ni siquiera para los que no usan paralelismo, como Regresión Logística o Naive Bayes), no es posible establecer una relación costo-beneficio real entre la ganancia en velocidad que ofrece la paralelización nativa y la mejora (o no) en las métricas de clasificación.

# **Explicación del procedimiento de limpieza del dataset.**

#### Descripción general del conjunto de datos

El conjunto de datos seleccionado y procesado corresponde a **PaySim1** (*Synthetic Financial Datasets For Fraud Detection*), disponible públicamente en la plataforma Kaggle. Este dataset fue generado a partir de un simulador basado en registros transaccionales reales extraídos de un servicio móvil de dinero financiero, desarrollado para investigar la detección de fraudes ante la escasez de datos financieros abiertos por motivos de confidencialidad.

El conjunto unificado se compone de **6,362,620 transacciones** y 11 variables (10 predictoras y 1 objetivo), cumpliendo holgadamente con el requisito mínimo de 1,000,000 de registros exigido en la rúbrica del curso. La variable objetivo, `isFraud`, es de naturaleza binaria (1 = transacción fraudulenta, 0 = transacción legítima). El conjunto presenta un **severo desbalance de clases**, donde únicamente **8,213 transacciones son fraudulentas**, representando aproximadamente el **0.1291%** del total de observaciones (frente a 6,354,407 transacciones legítimas, equivalentes al 99.8709%).

URL Dataset: <https://www.kaggle.com/datasets/ealaxi/paysim1>

URL Notebook de Limpieza en GitHub (Rama nuevo-dataset): <https://github.com/GimZtanCool/TB-Programacion_concurrente/blob/nuevo-dataset/limpieza_paysim_fraud.ipynb>

#### Integridad y diagnóstico de valores faltantes y duplicados

Se realizó un diagnóstico exhaustivo de integridad sobre las 6,362,620 filas:
* **Valores faltantes:** Se verificó que ninguna de las 11 columnas presenta valores nulos o vacíos (0% de nulos en la totalidad del conjunto). La naturaleza sintética y calibrada del simulador PaySim garantiza la completitud de los registros transaccionales.
* **Registros duplicados:** Se comprobó la inexistencia de filas completamente duplicadas (0 registros duplicados íntegros), garantizando que cada fila corresponde a un evento transaccional único en el tiempo.

#### Tratamiento de variables identificadoras y regla heurística

El dataset original incluye variables alfanuméricas identificadoras de cuentas y una regla fija tradicional:
* `nameOrig` (cuenta de origen) y `nameDest` (cuenta de destino): poseen cardinalidad masiva prácticamente individual (millones de identificadores únicos tipo `C1231006815` o `M1979787155`). Utilizar estas columnas directamente causaría un severo sobreajuste (*overfitting*) e incrementaría la dimensionalidad a millones de categorías en memoria.
* **Extracción semántica:** Se identificó que el prefijo del destinatario codifica el tipo de entidad receptora: `C` (Cliente particular / Customer) y `M` (Comercio / Merchant). Se extrajo la variable categórica binaria `dest_type` (`C` vs. `M`) y posteriormente se eliminaron las columnas de identificadores originales.
* **Evaluación y descarte de `isFlaggedFraud`:** Esta variable simula el sistema de control bancario tradicional basado en la regla fija de marcar transferencias individuales superiores a $200,000 USD. El análisis exploratorio reveló que de los 8,213 fraudes reales, la regla fija apenas marcó 16 transacciones (omitiendo 8,197 fraudes, con una tasa de falsos negativos del 99.80%). Por constituir una regla heurística ineficaz y redundante frente a los modelos de machine learning, fue eliminada de la matriz de entrenamiento.

#### Tratamiento de valores atípicos en montos (`amount`)

La variable `amount` presentó una marcada asimetría hacia la derecha (cola pesada). Mientras que la media general se sitúa en $179,861.90 USD, el monto promedio de las transacciones fraudulentas alcanza los **$1,467,967 USD**, frente a $178,197 USD en transacciones legítimas. Mediante el análisis de percentiles se determinó que el 99.99% de las observaciones se encuentra por debajo de **$9,615,000 USD**, con valores máximos extremos que superan los $92.4 millones de USD. Con el fin de evitar inestabilidades numéricas en el cálculo de gradientes durante el entrenamiento, se aplicó un criterio de recorte (*winsorizing*) acotando los montos superiores al percentil 99.99.

#### Ingeniería de características contables y temporales

A partir del análisis del dominio financiero y de la literatura especializada, se diseñaron cuatro variables predictoras de alto impacto:
1. **Discrepancia contable de origen (`errorBalanceOrig`):** En una operación normal, el saldo nuevo debe ser igual al saldo anterior menos el monto transferido ($\text{newbalanceOrig} = \text{oldbalanceOrg} - \text{amount}$). Se definió:
   $$\text{errorBalanceOrig} = \text{newbalanceOrig} + \text{amount} - \text{oldbalanceOrg}$$
   En transacciones legítimas esta discrepancia es nula o mínima, mientras que en transacciones fraudulentas los defraudadores suelen vaciar la cuenta a cero sin que coincida con el monto exacto.
2. **Discrepancia contable de destino (`errorBalanceDest`):** De forma análoga, se definió:
   $$\text{errorBalanceDest} = \text{oldbalanceDest} + \text{amount} - \text{newbalanceDest}$$
   Esta métrica captura anomalías donde la cuenta receptora no refleja el ingreso completo o presentaba saldo previo inexistente.
3. **Variables temporales cíclicas (`hora_del_dia` y `dia_de_la_semana`):** La variable original `step` representa horas acumuladas continuas ($1 \text{ paso} = 1 \text{ hora}$, abarcando 744 pasos correspondientes a un mes completo de 31 días). A partir de ella se derivaron:
   * $\text{hora\_del\_dia} = step \pmod{24} \in [0, 23]$
   * $\text{dia\_de\_la\_semana} = (step // 24) \pmod 7 \in [0, 6]$
   Estas variables permiten capturar patrones circadianos y estacionales propios de la actividad criminal digital.

#### Codificación de variables categóricas

Las variables categóricas nominales del dataset corresponden a `type` (con 5 modalidades: `CASH_IN`, `CASH_OUT`, `DEBIT`, `PAYMENT`, `TRANSFER`) y `dest_type` (`C`, `M`). El análisis exploratorio demostró que el fraude se concentra **única y exclusivamente** en los tipos `TRANSFER` (4,097 fraudes) y `CASH_OUT` (4,116 fraudes), siendo nulo en las restantes. Ambas variables fueron codificadas mediante *One-Hot Encoding* (`pd.get_dummies` con `drop_first=True`), transformándolas en columnas binarias con tipo numérico entero de bajo consumo de memoria (`int8`).

#### Balance de clases mediante ponderación

Dado que la clase minoritaria (fraude) representa únicamente el 0.1291% del dataset masivo de 6.36 millones de registros, el uso de técnicas de sobremuestreo sintético como SMOTE multiplicaría innecesariamente el volumen de datos en memoria a más de 12.7 millones de filas, colapsando los recursos durante el entrenamiento concurrente en Go previsto para la PC2. En su lugar, se implementó una estrategia de **ponderación de clases (*class weighting*)** en la función de pérdida:
$$W_j = \frac{N}{K \cdot N_j}$$
donde $N = 6,362,620$, $K = 2$, obteniéndose pesos de $W_0 \approx 0.5006$ para transacciones legítimas y $W_1 \approx 387.35$ para transacciones fraudulentas. Esto penaliza severamente los falsos negativos sin generar sobrecarga computacional.

#### División del conjunto de datos (Split temporal)

Debido a la naturaleza estrictamente secuencial y temporal de las transacciones financieras respecto de la variable `step`, la división del dataset en subconjuntos de entrenamiento y prueba se efectuó siguiendo un **criterio cronológico** (80% entrenamiento / 20% prueba):
* **Conjunto de entrenamiento (80%):** Primeras 5,090,096 transacciones cronológicas (pasos temporales anteriores), conteniendo 5,651 casos de fraude (0.1110%).
* **Conjunto de prueba (20%):** Últimas 1,272,524 transacciones cronológicas (pasos temporales más recientes), conteniendo 2,562 casos de fraude (0.2013%).
Este procedimiento previene la fuga de información temporal (*data leakage*), asegurando que los modelos evalúen transacciones futuras sin haber visto datos posteriores en su entrenamiento.

#### Estandarización

Las variables numéricas continuas y las derivadas de saldos y tiempos fueron estandarizadas con media cero y desviación estándar uno ($z = \frac{x - \mu}{\sigma}$). El escalador (`StandardScaler`) fue ajustado (*fit*) exclusivamente sobre el conjunto de entrenamiento (80%) y aplicado (*transform*) tanto al conjunto de entrenamiento como al de prueba, evitando cualquier filtración de información estadística hacia el conjunto de evaluación.

#### **Tabla 1**.

#### *Resumen del procedimiento de limpieza y preprocesamiento*

| **Etapa** | **Técnica aplicada** | **Justificación** |
| --- | --- | --- |
| **Diagnóstico de calidad** | Verificación de completitud y duplicados | Dataset sintético consistente: 0 nulos y 0 duplicados íntegros en 6,362,620 filas |
| **Tratamiento de identificadores** | Depuración de `nameOrig` y `nameDest`; extracción de prefijo `dest_type` | Prevenir sobreajuste por IDs individuales y capturar la distinción Cliente vs Comercio |
| **Depuración de regla estática** | Eliminación de `isFlaggedFraud` | Regla basada en umbral fijo (>200k) ineficaz (99.8% falsos negativos; 16 de 8,213 fraudes) |
| **Valores atípicos (`amount`)** | Recorte (*winsorizing*) al percentil 99.99 ($9,615,000 USD) | Controlar colas pesadas y asegurar la estabilidad numérica en el cálculo de gradientes |
| **Ingeniería contable** | Derivación de `errorBalanceOrig` y `errorBalanceDest` | Identificar discrepancias de saldo contable típicas en operaciones ilícitas |
| **Ingeniería temporal** | Derivación cíclica de `hora_del_dia` ($step \pmod{24}$) y `dia_de_la_semana` | Capturar estacionalidad diaria y semanal sin depender del contador lineal acumulado |
| **Variables categóricas** | One-Hot Encoding (`type`, `dest_type`) con `drop_first=True` | Representación binaria compacta (`int8`) compatible con modelos de aprendizaje automático |
| **Balance de clases** | Ponderación de clases en función de pérdida ($W_0 \approx 0.5006, W_1 \approx 387.35$) | Compensar el desbalance (0.13% fraude) sin duplicar memoria RAM ni sobrecargar Go |
| **División del conjunto** | Split cronológico 80/20 por variable `step` | Simular entorno de producción real y evitar fuga de información temporal (*data leakage*) |
| **Estandarización** | Z-score (`StandardScaler`) ajustado solo sobre entrenamiento | Escalamiento uniforme ($\mu=0, \sigma=1$) indispensable para convergencia y algoritmos métricos |

# **Trabajo en github, historial de commits de cada integrante del grupo.**

URL del repositorio: <https://github.com/GimZtanCool/TB-Programacion_concurrente>  
Rama oficial del entregable PC1: `nuevo-dataset` (<https://github.com/GimZtanCool/TB-Programacion_concurrente/tree/nuevo-dataset>)

El desarrollo colaborativo se estructuró siguiendo la metodología Git Flow:
* **Rama `main`:** Contiene la configuración base inicial y la documentación del proyecto.
* **Rama `nuevo-dataset`:** Rama de trabajo donde se integró el dataset masivo **PaySim1** (6,362,620 registros) gestionado mediante **Git LFS** (`paysim.csv`), junto con el notebook oficial de preprocesamiento (`limpieza_paysim_fraud.ipynb`), los scripts de partición y las muestras de validación (`data-limpia/`).

![](data:image/png;base64...)

# **Referencias.**

Salazar, A., Pérez, J. F., & Gallego, J. (2024). VigIA: Prioritizing public procurement oversight with machine learning models and risk indices. *Data & Policy, 6*, e75.<https://doi.org/10.1017/dap.2024.83>

Ukwu, A., Ghareeb, S., & Kazan, J. M. (2025). Financial fraud detection using machine learning models. In *2025 18th International Conference on Development in eSystem Engineering (DeSE)* (pp. 435-439). IEEE.<https://doi.org/10.1109/DeSE68208.2025.11367919>

Ghimire, A., & Amsaad, F. (2024). A parallel approach to enhance the performance of supervised machine learning realized in a multicore environment. *Machine Learning and Knowledge Extraction, 6*(3), 1840-1856.<https://doi.org/10.3390/make6030090>

Wu, X., Brazzle, P., & Cahoon, S. (2023). *Performance and energy consumption of parallel machine learning algorithms* (arXiv:2305.00798). arXiv.<https://arxiv.org/abs/2305.00798>