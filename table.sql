create database items;
create table purchases
(
    item_id int unique ,
    item_name varchar(225),
    item_quantity float,
    item_rate float,
    item_purchase_date date
);

insert into purchases values (1,'pencil',20,5.5,'2020-05-09');
insert into purchases values (1,'pen',10,5.5,'2020-04-04');
insert into purchases values (1,'rubber',5,5.5,'2019-01-01');