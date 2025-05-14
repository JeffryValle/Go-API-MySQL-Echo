use *******;

/* create table Clientes (
ClienteID int auto_increment primary key,
Nombre varchar(100) not null,
Apellido varchar(100) not null,
Edad int not null 
);

insert into Clientes values 
( 2,"Wakanda", "Forever", 30);

select * from Clientes;
*/

create table users (
    id int auto_increment primary key,
    name varchar(100) not null,
    email varchar(100) not null unique,
    password varchar(100) not null
);

select * from users;
