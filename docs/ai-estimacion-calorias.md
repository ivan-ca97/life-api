# ¿Qué tan precisa es la estimación de calorías por foto?

Si usás la función de análisis de comida por foto en Vitae, este artículo explica cómo funciona, qué tan precisa es, y cuándo conviene —o no— confiar en ella.

---

## Cómo funciona

Cuando fotografiás un plato, Vitae lo envía a un modelo de inteligencia artificial que analiza la imagen e intenta identificar los alimentos visibles, estimar las porciones, y calcular calorías y macronutrientes.

El modelo devuelve una estimación estructurada: calorías, proteínas, carbohidratos y grasas. Esos valores se cargan directamente como punto de partida, y vos podés ajustarlos antes de confirmar el registro.

---

## La precisión real — sin rodeos

### Calorías

El error promedio (MAPE) en estudios independientes para la estimación de calorías es de **alrededor del 35–40%**.

En la práctica, eso significa:

| Calorías reales | Rango probable de la estimación |
|----------------|-------------------------------|
| 400 kcal | 250 – 550 kcal |
| 700 kcal | 450 – 950 kcal |
| 1000 kcal | 630 – 1370 kcal |

Para una comida pequeña, el error puede ser de 150 kcal. Para un plato abundante, puede superar las 300 kcal.

### Macronutrientes

La estimación de proteínas, carbohidratos y grasas es considerablemente menos precisa que la de calorías:

| Nutriente | Error promedio |
|-----------|---------------|
| Proteína | 60–110% |
| Grasa | 40–90% |
| Carbohidratos | 50–70% |

Un plato con 30g de proteína podría ser estimado entre 12g y 60g. Estos valores son orientativos, no clínicamente confiables.

> Esta información proviene de un estudio peer-reviewed publicado en *Current Developments in Nutrition* (Oxford University Press, septiembre 2025), que evaluó los modelos de visión más avanzados disponibles sobre 52 fotografías estandarizadas con distintos tamaños de porción.

---

## Por qué es difícil — limitaciones estructurales

**El volumen y la densidad no son visibles en una foto.** Un bowl de arroz blanco y uno de arroz integral se ven idénticos. 200g de pollo a la plancha y 350g también. La IA no puede medir lo que no ve.

**Los platos mixtos multiplican el error.** En una cazuela, un guiso, o cualquier preparación donde los ingredientes no se distinguen claramente, cada componente individual acumula su propio margen de error.

**Las porciones grandes se subestiman sistemáticamente.** Los modelos tienden a subestimar cuándo hay mucha comida, especialmente si no hay referencia de escala en la imagen.

**La iluminación y el ángulo importan.** Una foto con mala luz o tomada desde muy arriba reduce significativamente la precisión.

---

## Cuándo es útil y cuándo no

### Útil

- Como punto de partida para no arrancar de cero cuando registrás
- Para alimentos simples y bien definidos: una fruta, un huevo, una porción de carne visible
- Para tener una idea general del orden de magnitud calórico de un plato nuevo
- Cuando el objetivo es llevar un registro aproximado y consistente a lo largo del tiempo

### No recomendado como única fuente

- Si seguís un plan nutricional con objetivos de macros específicos
- Para alimentación deportiva de alto rendimiento donde el margen importa
- Si tenés una condición clínica que requiere precisión (diabetes, trastornos alimentarios, etc.)
- Si querés confiar en los valores de proteína para evaluar una comida

---

## Credibilidad del registro

Cada alimento dentro de una comida tiene un nivel de confianza según cómo se registró su cantidad. Vitae lo muestra para que sepas qué tan confiables son los totales.

| Nivel | Cómo se registró | Error esperado |
|-------|-----------------|----------------|
| ⬜ Estimado por foto | La IA estimó la porción a partir de la imagen | ~35–40% |
| 🟡 Confirmado | Identificaste el alimento manualmente pero no pesaste | ~15–25% |
| 🟢 Pesado (cocinado) | Pesaste el ingrediente ya cocinado | ~10–15% |
| 🟢 Pesado en crudo | Pesaste el ingrediente antes de cocinarlo | ~8–12% |

El score de la comida es el del ingrediente menos preciso: si pesaste todo pero estimaste el aceite, el total hereda la incertidumbre del aceite.

Pesar en crudo es ligeramente más preciso que en cocinado porque las tablas nutricionales están referenciadas en crudo. Al cocinar, los alimentos pierden o absorben agua de manera variable según el método.

Para registrar el método de medición de cada ingrediente, usá el campo opcional dentro de cada ítem de la comida.

---

## Cómo sacar el mejor resultado

1. **Tomá la foto de frente, con buena luz.** Evitá contraluz y ángulos muy cenitales.
2. **Separar los componentes si es posible.** Si podés fotografiar el pollo aparte del arroz, la estimación mejora.
3. **Revisá siempre antes de confirmar.** Los valores pre-cargados son un borrador, no una medición.
4. **Para ingredientes conocidos, usá la base de alimentos.** Si sabés que comiste 150g de pechuga, cargar manualmente da un resultado más preciso que cualquier foto.
5. **La foto es útil para comidas en restaurantes o situaciones donde no tenés los datos exactos.** Ahí es donde aporta más valor.

---

## Sobre los modelos que usamos

El análisis de fotos en Vitae usa modelos de lenguaje de visión (VLMs) de última generación. Evaluamos distintas opciones en base a precisión documentada, costo y soporte de output estructurado.

Los mejores modelos disponibles hoy para esta tarea específica tienen un error de energía del ~35–40% en estudios controlados. Ningún modelo general de propósito amplio supera este umbral con robustez. Existen arquitecturas especializadas (sistemas RAG sobre bases de datos nutricionales como USDA FNDDS) que alcanzan mayor precisión, pero no están disponibles como APIs públicas.

Actualizamos periódicamente el modelo subyacente a medida que mejora el estado del arte.

---

## En resumen

La estimación por foto es una herramienta de conveniencia, no de precisión. Funciona mejor cuando se usa como punto de partida y se revisa antes de confirmar. Para un seguimiento nutricional serio, complementala siempre con tus propios registros o con datos de la base de alimentos.
