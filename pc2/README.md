# PC2: regresión logística ponderada para PaySim

La implementación de PC2 está en esta carpeta. El informe integrado se conserva localmente como `Informe_CC65-PC2-202620_Equipo.pdf` y no se publica en GitHub.

## Reproducción

Desde esta carpeta, con Go instalado y `paysim.csv` en la raíz del repositorio:

```powershell
go run . -input ..\paysim.csv -mode bench -repeats 30 -bench-epochs 1 -workers 1,2,4,8 -out mediciones.csv
go run . -input ..\paysim.csv -mode evaluate -epochs 5
```

El programa lee el CSV original completo, comprueba el orden cronológico, calcula el límite del percentil 99.99 de `amount`, ajusta la estandarización poblacional con el 80 % inicial y transforma el 20 % final reservado para prueba. Las 15 características y los pesos de clase siguen el preprocesamiento descrito en PC1.

Los tiempos del benchmark incluyen el cálculo del gradiente, la coordinación de workers, la reducción y la actualización del modelo. Excluyen la lectura del CSV y la estandarización. Las mediciones individuales se guardan en `mediciones.csv`; la consola muestra la media recortada al 10 %. El programa también informa la cantidad de CPU lógicas, la memoria asignada al heap de Go y el tiempo de CPU del proceso.

## Modelo Promela

`modelo_sincronizacion.pml` representa dos workers y un coordinador. Cada worker publica un resultado parcial privado; solo el coordinador consume los resultados y actualiza los parámetros compartidos. Si Spin está instalado, se pueden comprobar las aserciones y la ausencia de interbloqueos:

```powershell
spin -a modelo_sincronizacion.pml
gcc -O2 -o pan.exe pan.c
.\pan.exe -a
```
