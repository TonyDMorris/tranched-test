


CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX ON public.users (username);

CREATE TABLE IF NOT EXISTS public.assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES public.users(id) ON DELETE CASCADE,
    symbol TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

ALTER TABLE public.assets ADD CONSTRAINT amount CHECK (amount >= 0);

INSERT INTO public.users (username, password_hash) 
VALUES (
    'bob',
    '$2a$10$g6efNvur33Ya9lV1Fo3ClekXylbUgv2JEl.9.21vCFRPEoBopk6k.'
);

INSERT INTO public.assets (user_id, symbol, amount) 
VALUES (
    (SELECT id FROM public.users WHERE username = 'bob'),
    'EUR',
    10000.00
);

INSERT INTO public.assets (user_id, symbol, amount) 
VALUES (
    (SELECT id FROM public.users WHERE username = 'bob'),
    'USD',
    10000.00
);

INSERT INTO public.users (username, password_hash) 
VALUES (
    'tracy',
    '$2a$10$rK2Ia/SRslB0c7GTTsHKqewhFoccS/Q/189UHCkiT6HIMKg.lHgHi'
);

INSERT INTO public.assets (user_id, symbol, amount) 
VALUES (
    (SELECT id FROM public.users WHERE username = 'tracy'),
    'EUR',
    10000.00
);

INSERT INTO public.assets (user_id, symbol, amount) 
VALUES (
    (SELECT id FROM public.users WHERE username = 'tracy'),
    'USD',
    10000.00
);

