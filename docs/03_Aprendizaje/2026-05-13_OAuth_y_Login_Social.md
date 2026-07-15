# OAuth y login social

## Qué es realmente

Cuando la gente dice "login con Google", "login con Facebook" o "login con GitHub", normalmente está hablando de una combinación de:

- `OAuth 2.0` para delegar autorización
- `OpenID Connect (OIDC)` para identidad / login

En la práctica, para iniciar sesión de usuario final, lo correcto suele ser pensar en:

**OIDC sobre OAuth 2.0**

---

## Qué tan estándar es hoy

Sí, es completamente estándar en industria.

Especialmente para:

- apps web
- apps mobile
- productos B2C
- productos donde bajar fricción de registro importa mucho

También es común mezclar:

- login con email/password propio
- login social con Google/GitHub/Apple/etc.

---

## Qué tan complejo es

### Complejidad conceptual

Media.

No es “magia imposible”, pero tampoco es tan simple como un `POST /login` propio.

Porque hay que entender:

- redirecciones
- `state`
- `code`
- `callback URL`
- intercambio de tokens
- validación de identidad
- asociación de usuario externo con usuario interno

### Complejidad de implementación

Para PulseCity hoy diría:

- **baja-media** si usás una librería buena o un proveedor de auth
- **media-alta** si lo implementás bastante a mano

No es la parte más difícil del sistema, pero sí mete superficie nueva y varios edge cases.

---

## Qué se necesita

Como mínimo:

1. registrar tu app en Google / GitHub / etc.
2. obtener `client_id` y `client_secret`
3. definir una `redirect/callback URL`
4. agregar endpoints backend para:
   - iniciar login externo
   - recibir callback
   - intercambiar `authorization code` por tokens
5. validar identidad del proveedor
6. crear o encontrar un usuario interno en tu base
7. emitir tu propia sesión para PulseCity

O sea:

- el proveedor autentica al usuario
- **tu sistema sigue necesitando su propia sesión**

Google no reemplaza la sesión de PulseCity.

---

## Es pago

### En general

Normalmente **no pagás por usar OAuth/OIDC básico** con proveedores como Google o GitHub.

Lo que suele ser gratis:

- registrar la app
- usar login social estándar
- obtener identidad básica del usuario

### Dónde puede aparecer costo

- si usás un proveedor externo de auth tipo Auth0, Clerk, FusionAuth Cloud, etc.
- si superás límites free tier
- si querés features enterprise:
  - SSO corporativo
  - auditoría avanzada
  - MFA administrado
  - organizaciones / tenants complejos

Para un proyecto personal como PulseCity:

**Google login por sí solo no debería ser un costo importante**

---

## Qué ventajas tiene

- baja fricción de registro
- menos passwords propias que manejar
- mejor UX para muchos usuarios
- menos probabilidad de cuentas falsas por typo en email

---

## Qué desventajas tiene

- más complejidad que email/password simple
- dependés de redirecciones y configuración externa
- debugging un poco más molesto
- no elimina la necesidad de modelar usuarios y sesiones internas

Y además:

- no reemplaza ownership de partidas
- no reemplaza migraciones guest -> user
- no resuelve por sí solo auth de backend

---

## Qué recomendaría para PulseCity

Hoy, para el estado actual del proyecto, **no lo pondría todavía**.

Primero conviene cerrar bien:

- auth propia mínima
- sesiones
- ownership
- carga de partidas
- restore / logout

Después sí, agregar login con Google sería bastante razonable.

### Orden sensato

1. dejar sólido email/password propio
2. estabilizar modelo de usuario
3. recién después sumar Google como proveedor adicional

Eso hace que Google login sea un “nuevo método de entrada” y no una reescritura del sistema.

---

## Respuesta corta

### ¿Es estándar?

Sí.

### ¿Es complejo?

Complejidad media.

### ¿Qué se necesita?

Configurar proveedor externo + callback + validación + usuario interno + sesión propia.

### ¿Es pago?

Generalmente no para uso básico; depende más del proveedor de auth que del estándar OAuth en sí.

---

## Aplicado a PulseCity

Si más adelante querés, el camino más prolijo sería:

- mantener `email/password`
- sumar `login con Google`
- mapear Google user -> `users` existente
- reutilizar exactamente el mismo sistema de sesiones y ownership ya construido

Eso sería una extensión natural, no una refactorización traumática.
