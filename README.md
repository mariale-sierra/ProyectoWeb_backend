# 📦 Series Tracker - Backend

## App Corriendo
<img width="1867" height="883" alt="image" src="https://github.com/user-attachments/assets/853e815d-2460-42e3-937b-2c44c3e1d129" />
<img width="1753" height="494" alt="image" src="https://github.com/user-attachments/assets/38c7f9dc-6400-4497-9f86-3c5f2d7f5acd" />
<img width="1685" height="718" alt="image" src="https://github.com/user-attachments/assets/030440e6-a445-4abb-94d3-a3961b45620e" />



## 🚀 Descripción

Este es el backend del proyecto **Series Tracker**, una API REST construida en Go que permite gestionar series, progreso de episodios y ratings.

La API expone endpoints para:
- Listar series
- Crear nuevas series
- Actualizar progreso
- Eliminar series
- Asignar ratings

---

## 🛠️ Tecnologías utilizadas

- Go (net/http)
- SQLite3
- JSON (API REST)
- HTTP Server

---

## ▶️ Cómo ejecutar en el servidor
- Conectarse al servidor:
ssh -i id_gcp student@35.239.29.236
- Ir al backend:
cd ~/24405/ProyectoWeb/ProyectoWeb_backend
- Ejecutar el servidor:
nohup ./app &
- Abrir en el navegador:
http://35.239.29.236

## Reflexión

Para este proyecto se utilizó Go con su librería estándar net/http, lo cual permitió comprender mejor cómo funcionan los servidores HTTP sin depender de frameworks.

Trabajar sin frameworks hizo que conceptos como routing, manejo de requests, CORS y estructura de una API REST fueran mucho más claros. Aunque al inicio fue más complejo que usar herramientas modernas, el aprendizaje fue más profundo.

SQLite fue una buena elección por su simplicidad, aunque tiene limitaciones en entornos de producción real.

Los mayores retos fueron:

Manejo de CORS
Routing manual en Go
Deploy en servidor sin acceso a configuración de red

Definitivamente usaría Go nuevamente para APIs simples o proyectos donde se requiera control total sobre el backend.
