# Introducción

Este repositorio contiene el código básico que desarrollé para leer una cédula de ciudadanía o tarjeta de identidad utilizando el formato PDF-417. Está configurado para iniciar un servidor HTTP con endpoints, lo cual es útil si tienes una webapp y quieres escanear cédulas desde el cliente. Si solo necesitas el escáner, puedes eliminar el código del servidor y dejar abierto el puerto serial mientras la aplicación está en ejecución.

Debes utilizar un escáner de código de barras USB configurado para funcionar como dispositivo serial. Dependiendo de tu sistema operativo, es posible que necesites configurarlo para que emule un puerto serial.

Hay varias áreas que podrían mejorarse, como el manejo de las diferentes permutaciones de nombres (un apellido, un nombre, dos apellidos y un nombre, etc.). Aprecio cualquier PR que pueda ayudar a mejorar estas funcionalidades. Si algo no te funciona, no dudes en crear un issue con la cadena que está fallando para intentar darle soporte.

## Características

1. **Genera un ejecutable** que permite escanear un documento de identificación con formato PDF-417 y obtener la información del escaneo en formato JSON. Esta información puede ser consumida o consultada a través de endpoints.
2. **Inicia un servidor HTTP** en el puerto 1024, el cual expone dos endpoints.
3. **Abre un puerto serial** con un tiempo de espera (timeout) de 30 segundos para escanear una cédula.
4. **Formatea la entrada** del puerto serial en un objeto que contiene la información de la cédula.

## Endpoints

- **GET /getdata**: Este endpoint permite obtener los datos formateados de la cédula escaneada.
- **GET /cancel**: Este endpoint permite cancelar la operación de escaneo de la cédula.

# Instrucciones de Uso

Para compilar usa Go 1.22.3 o mayor.

1. **Instalación**:

   - Clona el repositorio en tu máquina local:
     ```bash
     git clone <URL_del_repositorio>
     ```
   - Navega hasta el directorio del proyecto:
     ```bash
     cd tu_repositorio
     ```

2. **Compilación**:

   - Si estás en Linux, puedes generar el ejecutable para Windows con:
     ```bash
     make
     ```
   - Para limpiar el build:
     ```bash
     make clean
     ```
   - En Windows, haz un build normal:
     ```cmd
     go build -o barcode_scanner
     ```

3. **Uso de Endpoints**:
   - Para obtener los datos de la cédula:
     ```http
     GET http://localhost:1024/getdata
     ```
   - Para cancelar el escaneo de la cédula:
     ```http
     GET http://localhost:1024/cancel
     ```

## Notas

- Asegúrate de que el puerto 1024 esté disponible y no sea utilizado por otra aplicación.
- El tiempo de espera (timeout) del puerto serial es de 30 segundos.
- Los datos de la cédula serán formateados en un objeto JSON.

## Contribuir

Si tienes ideas para mejorar el código, añadir nuevas funcionalidades, o has encontrado algún problema, no dudes en abrir un issue o hacer un pull request.

Para contribuir:

1. Haz un fork de este repositorio.
2. Crea una nueva rama:

    ```bash
    git checkout -b mi-nueva-rama
    ```

3. Realiza tus cambios y haz commit:

    ```bash
    git commit -m "Descripción de los cambios"
    ```

4. Sube tu rama:

    ```bash
    git push origin mi-nueva-rama
    ```

5. Abre un pull request.

## Créditos

Agradecimientos a [Eitol](https://github.com/Eitol/colombian-cedula-reader) por proporcionar la lista de ciudades con el código DIVIPOL y por la información sobre cómo se organiza la cadena leída por el escáner. Esta implementación se basa en su trabajo, pero en Go. La forma en que extraigo la información de las cadenas es diferente: mientras que Eitol utiliza rangos de bytes, yo decidí usar un puntero, guiándome por ciertos caracteres que indican el inicio y el final de los datos a extraer.

