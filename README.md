# Shop Basket Project

## Development Setup

Copy the the contents of `.env.example` in to a newly created `.env` file and adjust the values.

run `docker-compose up -d` for postgresql

## Running the Server
Navigate to the project directory and run the following command to start the server:

`go run main.go`

The server will start listening on port `3005` by default. You can change the port in the `main.go` file if needed.

### Test Endpoints
* **Register**

    ``` bash
    curl -d "username=`your_username`" -d "password=`your_password`" http://localhost:3005/register
    ```

* **Login**

    ``` bash
    curl -d "username=me" -d "password=1234" http://localhost:3005/login
    ```

* **Get Your Baskets**
    ``` bash
    curl -H "Authorization: Bearer `your_token`" http://localhost:3005/basket
    ```

* **Get a Basket with ID**
    ``` bash
    curl -H "Authorization: Bearer `your_token`" http://localhost:3005/basket/:id
    ```

* **Create a Basket**
    ``` bash
    curl -H "Authorization: Bearer `your_token`" -d "data=`json_data`"  -d "item=`json_item`" http://localhost:3005/basket/:id
    ```
* **Modify a Basket**
    ``` bash
    curl -H "Authorization: Bearer `your_token`" -X PATCH -d "data=`json_data`"  -d "item=`json_item`" http://localhost:3005/basket/:id
    ```

* **Delete a Basket**
    ``` bash
    curl -H "Authorization: Bearer `your_token`" -X DELETE http://localhost:3005/basket/:id
    ```

---
### author: @mpnri

