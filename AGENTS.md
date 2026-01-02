# AGENTS.md

Proposito
- Este documento guia a Codex para trabajar en este repo.
- Prioriza cambios pequenos, claros y faciles de revisar.

Contexto del repo
- Implementa un shell POSIX basico en Go (CodeCrafters).
- Entrada principal: `app/main.go`.
- Logica del shell en `internal/shell/` y comandos built-in en `internal/command/`.
- Parser/lexer en `internal/parser/` y `internal/lexer/`.

Buenas practicas Go
- Formato: ejecutar `gofmt` en todo archivo editado.
- Imports: mantener orden estandar; evitar imports sin uso.
- Errores: preferir `if err != nil { ... }` y devolver contexto con `fmt.Errorf("...")`.
- Interfaces: pequenas y enfocadas; recibir interfaces, devolver structs.
- Context: propagar `context.Context` en operaciones largas o cancelables.
- Concurrencia: usar goroutines solo si es necesario; proteger estado compartido.
- Logging: usar `fmt` para salida de usuario y `log` solo si aporta valor; no mezclar estilos.
- Panics: evitar; solo en condiciones irreparables.
- Estado global: evitar variables globales mutables.

Estilo y estructura
- Nombres de paquetes en minusculas y sin guiones.
- Funciones pequenas; extraer helpers si una funcion crece demasiado.
- Comentarios: solo cuando el codigo no sea obvio.
- Usar ASCII en archivos nuevos a menos que el repo ya use Unicode.

Testing y validacion
- Preferir tests de tabla (`[]struct{...}`) cuando aplique.
- Si se agregan tests, ejecutar `go test ./...`.
- Si no hay tests, validar comportamiento con casos manuales basicos.

Reglas especificas del shell
- Tokenizacion: mantener compatibilidad con comillas simples/dobles y escapes.
- Redirecciones: respetar `>`, `>>`, `2>`, `2>>` y restaurar stdout/stderr.
- Builtins: asegurar que la salida sea correcta y consistente con el shell.
- Ejecutables externos: usar `exec.LookPath` y preservar `Args[0]`.
- Interfaz: el prompt es `$ ` y se debe mantener el flujo REPL.

Flujo de trabajo sugerido
- Leer archivos relevantes antes de tocar codigo.
- Cambiar lo minimo posible para cumplir el objetivo.
- Revisar impacto en otras partes del shell (lexer/parser/runtime).
- Aplicar `gofmt` y reportar pasos de verificacion.
