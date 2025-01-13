## Build
A custom `migration.sh` shell script instead of a `Makefile` is used to build the project.  The results from both
are identical, just that a shell script is more powerful in what it can do and generally can be implemented more cleanly than an equivalent `Makefile`.  Of course, a `Makefile` can delegate to a shell script for its targets, but that
implies we already need a shell script.

### Migration
###### Create a     Blank Migration file
Pass the name of the migration file without any pace 
```bash
sh ./migration.sh migration create_user_table
```

###### Migrate
##### 1: Migration up
To create all database tables
```bash
sh ./migration.sh migrate up
```

##### 2: Migration down
To remove all database tables
```bash
sh ./migration.sh migrate down
```

##### 3: Seeder seed
```bash
sh ./migration.sh seeder seed
```

##### 4: Migration force
###### version of migration
```bash
sh ./migration.sh migrate force 1
```

