# Verificación formal del protocolo de reducción de gradientes

## Correspondencia con Go

`modelo_tp.pml` abstrae una época de `parallelGradient` seguida por `update`. El cálculo numérico se representa como una tarea terminada, sin valores flotantes ni millones de filas. Se comprueba coordinación, exclusividad de escritura y finalización.

| Elemento Promela | Correspondencia Go |
| --- | --- |
| Producer y jobs | Envío de W intervalos disjuntos y close(jobs) |
| Worker y ready[id] | accumulate produce un parcial privado por trabajo |
| results | Canal de gradientes con capacidad W |
| workers_done | Contador sincronizado de WaitGroup.Done |
| Closer y results_closed | wg.Wait seguido por close(results) |
| Coordinator y consumed[id] | Recepción y reducción de todos los parciales |
| writers y updates | Actualización exclusiva por el coordinador |

Los identificadores corresponden a trabajos, no a workers: la cola puede entregar varias tareas a uno. El coordinador espera el cierre de trabajos antes de recibir, como Go. La recepción sigue permitida después del cierre mientras queden mensajes, igual que range en un canal cerrado.

Solo el arranque de procesos está dentro de atomic. El protocolo y la escritura pueden intercalarse. `workers_done++` abstrae el contador sincronizado de WaitGroup; no propone un incremento desprotegido en Go. Los arrays ready y consumed son instrumentación, no gradientes compartidos de la implementación.

## Propiedades

Seguridad: Pan se compila con `-DSAFETY -DNOCLAIM` para comprobar aserciones y estados finales inválidos. Cada resultado debe existir, consumirse una sola vez y pertenecer a un trabajo válido. Antes de actualizar deben terminar W workers y recibirse W resultados, todos consumidos. writers debe ser uno durante la actualización y cero antes; updates debe ser exactamente uno.

Vivacidad: ejecución separada con `pan -a` y `ltl completion { <> finished }`. Toda ejecución debe alcanzar la finalización. No se habilita fairness; el modelo correcto es finito y no contiene ciclos de espera activos infinitos. No se usa BITSTATE: se almacena el espacio de estados con reducción de orden parcial. Interpretación conforme a [Pan](https://spinroot.com/spin/Man/Pan.html) y [LTL](https://spinroot.com/spin/Man/ltl.html).

## Resultados

Los conteos exactos están en [spin_resumen.csv](evidencias/spin_resumen.csv). Las seis búsquedas correctas con W=1,2,3 terminaron con cero errores, sin alcanzar el límite de profundidad y sin búsquedas incompletas. Los logs conservan estados, transiciones y profundidad. El entorno fue Spin 6.5.2 y GCC 12.2.0 en Debian bajo Docker/WSL2.

Los logs de seguridad muestran `invalid end states +`; los de vivacidad verifican el claim completion. Los estados no alcanzados del never claim no indican un fallo: corresponden al autómata negado cuya ejecución aceptante no se encontró. Los estados de los procesos del protocolo fueron alcanzados en la exploración reportada.

## Controles defectuosos

BAD_WRITER permite a los workers escribir: con W=2 se obtiene `assertion violated (writers==1)`. BAD_WAIT omite la finalización de un worker: el cierre queda bloqueado y aparece `invalid end state`. Se conserva una traza .trail y su reproducción textual para cada control.

Estos controles se detienen en el primer contraejemplo y muestran `Search not completed`. Es lo esperado y no se presenta como búsqueda completa del modelo correcto.

## Limitaciones

La exclusión mutua se obtiene por propiedad exclusiva; no se agregó un mutex innecesario. Esto no prueba exactitud numérica, manejo del CSV, fallos del sistema operativo, cancelación ni cualquier cantidad de workers. No es una traducción automática certificada de Go. Las pruebas Go complementan el modelo.

Se representa una época, W trabajos y capacidad W en ambos canales. Go crea canales y workers nuevos en cada época. Un pool persistente, buffers distintos, tareas dinámicas o errores de workers requieren ampliar el modelo y repetir las pruebas.
