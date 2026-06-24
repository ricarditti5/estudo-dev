--
-- PostgreSQL database dump
--

\restrict MK0a4J7u4TbDizPIR0v4VCsXgeZ0rBnRPWakHvCYjwpev0YOKN5cXYUdzufpGxe

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

-- Started on 2026-06-24 08:55:35

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 860 (class 1247 OID 25015)
-- Name: role_users; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.role_users AS ENUM (
    'admin',
    'user'
);


ALTER TYPE public.role_users OWNER TO postgres;

--
-- TOC entry 857 (class 1247 OID 24929)
-- Name: task_status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.task_status_enum AS ENUM (
    'pending',
    'in-progress',
    'done'
);


ALTER TYPE public.task_status_enum OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 219 (class 1259 OID 24915)
-- Name: tasks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tasks (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    description character varying(255) NOT NULL,
    status public.task_status_enum DEFAULT 'pending'::public.task_status_enum,
    tags text[],
    user_id uuid
);


ALTER TABLE public.tasks OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 24926)
-- Name: tasks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.tasks ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.tasks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 221 (class 1259 OID 25031)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nome character varying(100) NOT NULL,
    role public.role_users DEFAULT 'user'::public.role_users NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 5022 (class 0 OID 24915)
-- Dependencies: 219
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tasks (id, title, description, status, tags, user_id) FROM stdin;
109	Tester 1	New constraints tester	in-progress	{tst,newTst}	15d1cefa-ef63-4646-a451-af37e5fd7df7
110	Tester 2	Verificação de fluxo de dados	in-progress	{fluxo,backend}	15d1cefa-ef63-4646-a451-af37e5fd7df7
112	Tester 4	Correção de bug de autenticação	in-progress	{bug,seguranca}	15d1cefa-ef63-4646-a451-af37e5fd7df7
114	Tester 6	Implementação do layout responsivo	in-progress	{css,design}	15d1cefa-ef63-4646-a451-af37e5fd7df7
116	Tester 8	Testes unitários de serviços de email	in-progress	{testes,email}	15d1cefa-ef63-4646-a451-af37e5fd7df7
118	Tester 10	Integração com gateway de pagamento	in-progress	{pagamento,stripe}	15d1cefa-ef63-4646-a451-af37e5fd7df7
120	Tester 12	Ajuste de permissões do utilizador	in-progress	{acl,auth}	15d1cefa-ef63-4646-a451-af37e5fd7df7
122	Tester 14	Limpeza de logs antigos do servidor	in-progress	{manutencao,logs}	15d1cefa-ef63-4646-a451-af37e5fd7df7
124	Tester 16	Monitorização de consumo de memória	in-progress	{perf,servidor}	15d1cefa-ef63-4646-a451-af37e5fd7df7
126	Tester 18	Tradução do painel para Inglês	in-progress	{i18n,frontend}	15d1cefa-ef63-4646-a451-af37e5fd7df7
128	Tester 20	Deploy do ambiente de staging	in-progress	{deploy,staging}	15d1cefa-ef63-4646-a451-af37e5fd7df7
\.


--
-- TOC entry 5024 (class 0 OID 25031)
-- Dependencies: 221
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, nome, role) FROM stdin;
15d1cefa-ef63-4646-a451-af37e5fd7df7	Ricardo	admin
692920e0-8c85-4feb-a8c8-1cdd5b2a38fb	Alexandre	user
\.


--
-- TOC entry 5030 (class 0 OID 0)
-- Dependencies: 220
-- Name: tasks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tasks_id_seq', 128, true);


--
-- TOC entry 4871 (class 2606 OID 24922)
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- TOC entry 4873 (class 2606 OID 25040)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4869 (class 1259 OID 25005)
-- Name: idx_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_status ON public.tasks USING btree (status);


--
-- TOC entry 4874 (class 2606 OID 25043)
-- Name: tasks fk_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


-- Completed on 2026-06-24 08:55:35

--
-- PostgreSQL database dump complete
--

\unrestrict MK0a4J7u4TbDizPIR0v4VCsXgeZ0rBnRPWakHvCYjwpev0YOKN5cXYUdzufpGxe

