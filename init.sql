CREATE TABLE public.movies
(
    id integer NOT NULL,
    name character varying NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.movies
    OWNER to postgres;
