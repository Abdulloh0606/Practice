create database minitrello;
create table users(
	id serial primary key,
	name varchar(50) not null,
	email varchar(100) unique not null,
	password_hash text not null,
	role varchar(20) default 'user',
	created_at timestamp default current_timestamp
);

create table projects(
	id serial primary key,
	name varchar(25) not null,
	created_by int references users(id) on delete cascade,
	created_at timestamp default current_timestamp
);

create table tasks(
	id serial primary key,
	name varchar(50) not null,
	description text,
	comment text,
	status varchar(20) default 'TO DO',
	project_id int references projects(id),
	created_at timestamp default current_timestamp
);
CREATE TABLE project_members (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'user', 
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (project_id, user_id)      
);
alter table tasks add column deadline timestamp; 
alter table tasks add column assigned_to int references users(id); 
