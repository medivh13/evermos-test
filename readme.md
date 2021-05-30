# evermos backend test

# problem
We are members of the engineering team of an online store. When we look at ratings for our online store application, we received the following
facts:
- Customers were able to put items in their cart, check out, and then pay. After several days, many of our customers received calls from
our Customer Service department stating that their orders have been canceled due to stock unavailability.
- These bad reviews generally come within a week after our 12.12 event, in which we held a large flash sale and set up other major
discounts to promote our store.

After checking in with our Customer Service and Order Processing departments, we received the following additional facts:
- Our inventory quantities are often misreported, and some items even go as far as having a negative inventory quantity.
The misreported items are those that performed very well on our 12.12 event.
- Because of these misreported inventory quantities, the Order Processing department was unable to fulfill a lot of orders, and thus
requested help from our Customer Service department to call our customers and notify them that we have had to cancel their orders.

# My Solution

- the quantity of product's stock have to be updated when any customer do CHECKOUT, so the amount of remaining stock has changed. This will prevent any customer from taking any product beyond the remaining stock.
- the quantity of product's stock have to be updated when any customer do CANCEL the order, so the remaining stock will increase.
- the quantity of product's stock also have to be update when any order get expired due to the due_date (in this case, I set it 1 day), so the remaining stock will also increase. this is supposed to be done by CronJob, by Cron hiting the endpoints I have provided in serveral interval of time. Considering that in this test I couldn't create a cronjob on your server, so I only created the endpoint. which in practice this endpoint will be executed by CronJob.
- order process is have to be do in DB transaction, so in case any problem or error when inserting order data and or order detail data, the DB can do the ROLLBACK

# PS
- the test number 2 (treasure hunt) I made it on a separate repo. you can see in my reply email.

## Prerequisite:

1. This project requires Golang to run, view installation instruction here https://golang.org/doc/install
2. This project using postgresql. the table is in public scheme. create the table like sql example in pkg/database/tabledump.sql

## How to run:

This project depedency is managed by go mod, therefore you can put it anywhere other than GO-PATH directory.
steps to run:

- configure the DB in config/config-dev.yaml
- cd to project root directory 
- then do "go run main.go"
- all route is in transport/http/handlers
- all endpoints are protected by signature. since is a development env, I return the valid signature if you input the wrong one, and you can use the valid one. the signature formula is in all DTO validation (pkg/dto/..) in case you wanna try to make a valid signature by your self.

## Example

## Example

I give the curl example of all endpoint, but if you mind to use the curl in terminal or in Postman, you can just import this Postman collection via this link :

- postman import link : https://www.getpostman.com/collections/b1de44e3e6dbb9fea34e

# Curl Example

- ping : 
curl --location --request GET 'http://localhost:9009/api/evermos-test/ping'

- register :
curl --location --request POST 'http://localhost:9009/api/evermos-test/register' \
--header 'signature: 9c99b2f7b59409c87e1da968307c4dd747528cf233ae6078266eddeb3793e1a9' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email" : "jody.almaida@gmail.com",
    "password" : "password123"
}'

- login :
curl --location --request POST 'http://localhost:9009/api/evermos-test/login' \
--header 'signature: 9c99b2f7b59409c87e1da968307c4dd747528cf233ae6078266eddeb3793e1a9' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email" : "jody.almaida@gmail.com",
    "password" : "password123"
}'

- get all products :
curl --location --request GET 'localhost:9009/api/evermos-test/products' \
--header 'signature: a64034d779ada43e88ca97005b56cc95533b4be8959f90db96deb649e920e4d9' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjIzNTQ5NDQsInVzZXJfaWQiOjF9.pVjSGne8OyaiyqQox_HTqaCNaUO2onWxRWD9zQuTH64' \
--data-raw ''

- order :
curl --location --request POST 'localhost:9009/api/evermos-test/order' \
--header 'signature: a64034d779ada43e88ca97005b56cc95533b4be8959f90db96deb649e920e4d9' \
--header 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjIzNTQ5NDQsInVzZXJfaWQiOjF9.pVjSGne8OyaiyqQox_HTqaCNaUO2onWxRWD9zQuTH64' \
--header 'Content-Type: application/json' \
--data-raw '{
    "product_id" : [3,4],
    "quantity" : [2,3]
}'

- cancel order by id :
curl --location --request POST 'localhost:9009/api/evermos-test/cancel' \
--header 'signature: a64034d779ada43e88ca97005b56cc95533b4be8959f90db96deb649e920e4d9' \
--header 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjIzNTQ5NDQsInVzZXJfaWQiOjF9.pVjSGne8OyaiyqQox_HTqaCNaUO2onWxRWD9zQuTH64' \
--header 'Content-Type: application/json' \
--data-raw '{
    "order_id" : 2
}'

- cancel all expired checkout (order) :
curl --location --request POST 'localhost:9009/api/evermos-test/expired' \
--data-raw ''


