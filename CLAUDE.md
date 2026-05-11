# Общие правила

- Не спрашивай разрешение на создание/изменение файлов с раширением `.go`.
- Always use Context7 when I need library/API documentation, code generation, setup or configuration steps without me having to explicitly ask.
- use skill: golang-design-patterns

# References

TODO: если будут ссылки на документацию, то необходимо дополнить


# Critical Constraint

Never use log.Fatal or os.Exit outside of main().
Return errors. We instrument error rates via middleware
and os.Exit bypasses it entirely.

