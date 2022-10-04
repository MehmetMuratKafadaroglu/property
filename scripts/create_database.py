import psycopg2
DATABASE = "property"

def delete_tables():
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("""DROP TABLE Users CASCADE;
                DROP TABLE Properties CASCADE;
                DROP TABLE TemporaryKeys CASCADE;
                DROP TABLE Saved CASCADE;
                DROP SEQUENCE propertyimages_sequence;
                DROP TABLE PropertyImages CASCADE;""")
    conn.commit()
    conn.close()

def create_tables(only_create=False):
        try:
                if not only_create:
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
                ID BIGINT PRIMARY KEY,
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
                propertyType VARCHAR(20),
                isPublished BOOLEAN);

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
                userID INTEGER REFERENCES Users(ID) ON DELETE CASCADE,
                key CHAR(128));
                CREATE INDEX temporarykeys_userid_index ON TemporaryKeys USING HASH(userID);   
                """)
        print("Creating Saved")
        
        cur.execute("""
                CREATE TABLE Saved(
                ID SERIAL PRIMARY KEY,
                propertyID INTEGER REFERENCES Properties(ID) ON DELETE CASCADE,
                userID INTEGER REFERENCES Users(ID) ON DELETE CASCADE );
                CREATE INDEX saved_userid_index ON Saved(userID) ;
        """)
        print("Creating PropertyImages")
        cur.execute("""
                CREATE TABLE PropertyImages(
                ID BIGINT PRIMARY KEY,
                propertyID INTEGER REFERENCES Properties(ID) ON DELETE CASCADE,
                fileName VARCHAR(255)
                );
        """)
        conn.commit()
        cur.execute("CREATE SEQUENCE propertyimages_sequence AS BIGINT OWNED BY PropertyImages.ID;")
        cur.execute("CREATE SEQUENCE propertiesid_sequence AS BIGINT OWNED BY Properties.ID;")
        conn.commit()
        cur.close()
        conn.close()

if __name__ == "__main__":
        create_tables()