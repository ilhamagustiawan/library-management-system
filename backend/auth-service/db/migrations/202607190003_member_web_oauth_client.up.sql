-- Fail on a target-ID collision instead of deleting either client.
UPDATE oauth_clients
SET id = 'member-nextjs-web',
    name = 'Library Management System for Members'
WHERE id = 'nextjs';
