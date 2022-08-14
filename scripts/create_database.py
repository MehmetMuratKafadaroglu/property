import psycopg2
DATABASE = "property"

def delete_tables():
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    print("Deleting Users")
    cur.execute("DROP TABLE Users CASCADE")
    print("Deleting Properties")
    cur.execute("DROP TABLE Properties CASCADE") 
    print("Deleting TemporaryKeys")
    cur.execute("DROP TABLE TemporaryKeys CASCADE") 
    print("Deleting TemporaryKeys")
    cur.execute("DROP TABLE Saved CASCADE") 
    print("Deleting Saved")
    conn.commit()
    conn.close()

def create_tables():
        try:
                delete_tables()
        except:
                pass
        conn = psycopg2.connect(host = "localhost", 
                                dbname=DATABASE, 
                                user="postgres",
                                password="630991")
        cur = conn.cursor()
        print("Creating Properties Users")
        cur.execute("""
                CREATE TABLE Users(
                ID SERIAL PRIMARY KEY,
                email VARCHAR(255),
                companyName VARCHAR(255),
                isAgent BOOLEAN,
                phoneNumber bigint,
                password VARCHAR(255),
                isMailVerified bool) 
                """)
        print("Creating Properties Table")
        cur.execute("""
                CREATE TABLE Properties(
                ID SERIAL PRIMARY KEY,
                price INTEGER,
                isForSale BOOLEAN,
                numberOfRooms INTEGER,
                location VARCHAR(255),
                address  VARCHAR(1024),
                internalArea INTEGER,
                title VARCHAR(255),
                description VARCHAR(10000),
                publishDate DATE,
                authorID INTEGER REFERENCES Users(ID),
                orienter VARCHAR(1024),
                propertyType CHAR(20));

                CREATE INDEX properties_price_index ON Properties(price);
                CREATE INDEX properties_numberofrooms_index ON Properties(numberOfRooms);
                CREATE INDEX properties_location_index ON Properties USING HASH(authorID);   
                CREATE INDEX properties_internalarea_index ON Properties(internalArea);
                CREATE INDEX properties_publishdate_index ON Properties(publishDate);

                """)
        print("Creating TemporaryKeys")
        cur.execute("""
                CREATE TABLE TemporaryKeys(
                ID SERIAL PRIMARY KEY,
                userID INTEGER REFERENCES Users(ID),
                key CHAR(128));
                CREATE INDEX temporarykeys_userid_index ON TemporaryKeys USING HASH(userID);   
                """)
        print("Creating Saved")
        cur.execute("""
                CREATE TABLE Saved(
                ID SERIAL PRIMARY KEY,
                propertyID INTEGER REFERENCES Properties(ID),
                userID INTEGER REFERENCES Users(ID) );
                CREATE INDEX saved_userid_index ON Saved(userID);
        """)
        conn.commit()
        cur.close()
        conn.close()
