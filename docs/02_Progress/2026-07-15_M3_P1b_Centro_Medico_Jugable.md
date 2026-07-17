# M3.P1b — Centro Medico jugable

## Objetivo de la sesion

Hacer operable el flujo medico desde el frontend y cerrar el segundo hueco de gameplay detectado en la pausa de analisis de M3, sin adelantar los paneles completos de `M3.22`.

## Decision de UX

Centro Medico vive en una pagina propia, igual que Trade Center. El Command Center resume y navega; no absorbe una mecanica que requiere revisar varios casos, comparar riesgo y tomar decisiones con seguimiento.

La pantalla usa una direccion visual clinica y operativa dentro del design system: densidad media-alta, superficies oscuras, tipografia editorial para nombres/titulares y colores semanticos para diferenciar protocolo, precaucion y riesgo.

## Implementacion

- nueva ruta `/franchise/medical`
- nueva pagina `MedicalCenterPage`
- acceso Command Center → Centro Medico y retorno explicito
- resumen de jugadores disponibles, lesiones activas y casos graves
- casos activos con nombre, posicion, rating, severidad, dias estimados y retorno esperado
- cuatro decisiones operables por lesion:
  - seguir protocolo (`rest`)
  - carga reducida al volver (`reduce_minutes`)
  - ignorar recomendacion (`ignore_doctor`)
  - forzar alta anticipada (`force_return`)
- estados de espera de snapshot, roster sano, envio, error y confirmacion
- hook de dominio `useMedicalOperations`
- estado de requests indexado por `injury_id`
- CSS propio y responsive en `components/medical/medicalCenter.css`

## Ownership y flujo

La pagina deriva disponibilidad desde `roster.patch`. El frontend no diagnostica, no calcula riesgo y no cambia el roster.

Al decidir, `useMedicalOperations` publica al endpoint existente del gateway. La respuesta HTTP confirma que la decision fue registrada. La disponibilidad no cambia de forma optimista: `team-service` conserva la autoridad y la pagina espera el delta posterior. Esto importa especialmente porque `rest`, `reduce_minutes` e `ignore_doctor` persisten una decision sin producir un cambio inmediato de disponibilidad, mientras que `force_return` si puede generar un `roster.patch` de alta.

## Verificacion

- `npm run build --prefix frontend`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `make build`

Todos finalizaron correctamente.

## Limites conocidos

- El roster contractual completo sigue sin rehidratarse por REST; esta pagina usa los deltas vivos disponibles, como el Trade Center de `M3.P1a`.
- La confirmacion local de una decision medica no se rehidrata al recargar. La persistencia backend existe; su lectura agregada queda para los paneles/inbox de `M3.22` y `M3.23`.
- El playtest manual no forma parte de este corte de implementacion.

## Pendiente siguiente

`M3.P1c`: jugar una temporada a mano, proponer al menos tres trades, responder las decisiones medicas que aparezcan y documentar game feel como input para `M3.13`.
