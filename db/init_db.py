import sqlite3

drop_snippets_stmt = '''DROP TABLE IF EXISTS snippets;'''
create_snippet_stmt = '''CREATE TABLE snippets (
								id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
								title TEXT NOT NULL,
								content TEXT NOT NULL,
								created INTEGER NOT NULL,
								expires INTEGER NOT NULL
							);'''
index_snippets_created = '''CREATE INDEX idx_snippets_created ON snippets(created);'''

drop_sessions_stmt = '''DROP TABLE IF EXISTS sessions;'''
create_sessions_stmt = '''CREATE TABLE sessions (
								token TEXT PRIMARY KEY,
								data BLOB NOT NULL,
								expiry REAL NOT NULL
							);'''
index_sessions_expiry = '''CREATE INDEX sessions_expiry_idx ON sessions(expiry);'''

def init_db():
    conn = sqlite3.connect('./db/snippetbox.db')
    c = conn.cursor()

    c.execute(drop_snippets_stmt)
    c.execute(create_snippet_stmt)

    c.execute(drop_sessions_stmt)
    c.execute(create_sessions_stmt)

    conn.commit()
    conn.close()

if __name__ == '__main__':
    init_db()