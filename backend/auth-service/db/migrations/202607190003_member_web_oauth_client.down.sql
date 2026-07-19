-- Fail on a legacy-ID collision instead of deleting either client.
UPDATE oauth_clients
SET id = 'nextjs',
    name = 'Library Management System Next.js'
WHERE id = 'member-nextjs-web';
