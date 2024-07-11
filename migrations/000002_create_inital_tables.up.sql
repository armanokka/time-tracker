create table "user"
(
    id              bigserial
        primary key,
    email           text   not null
        constraint uni_user_email
            unique,
    password        text   not null,
    name            text   not null,
    surname         text   not null,
    patronymic      text,
    address         text,
    admin bool not null default false
);

-- alter table user
--     owner to effective_mobile;



create table project
(
    id          bigserial
        primary key,
    name        text   not null,
    description text,
    creator_id  bigint not null
        constraint fk_project_user
            references "user"
            on update cascade on delete cascade
);

-- alter table project
--     owner to effective_mobile;

create unique index name_creator_id_idx
    on project (name, creator_id);

create table project_participant
(
    user_id    bigint not null
        constraint fk_project_participant_user
            references "user"
            on update cascade on delete cascade,
    project_id bigint not null
        constraint fk_project_participant_project
            references project
            on update cascade on delete cascade
);

-- alter table project_participant
--     owner to effective_mobile;

create unique index user_id_project_id_idx
    on project_participant (user_id, project_id);


create table task
(
    id          bigserial
        primary key,
    name        text                                               not null,
    description text,
    project_id  bigint                                             not null
        constraint fk_task_project
            references project
            on update cascade on delete cascade,
    finished boolean default false not null
);

-- alter table task
--     owner to effective_mobile;


create table task_participant
(
    task_id bigint not null
        constraint fk_task_participant_task
            references task
            on update cascade on delete cascade,
    user_id bigint not null
        constraint fk_task_participant_user
            references "user"
            on update cascade on delete cascade
);

-- alter table task_participant
--     owner to effective_mobile;

create unique index user_id_task_id_idx
    on task_participant (task_id, user_id);

create table time_entry
(
    id         bigserial
        primary key,
    task_id    bigint not null
        constraint fk_time_entry_task
            references task
            on update cascade on delete cascade,
    user_id    bigint not null
        constraint fk_time_entry_user
            references "user"
            on update cascade on delete cascade,
    started_at timestamp with time zone default CURRENT_TIMESTAMP,
    ended_at   timestamp with time zone
);

alter table time_entry
    add constraint check_time
        check (started_at <= ended_at);



-- alter table time_entry
--     owner to effective_mobile;

