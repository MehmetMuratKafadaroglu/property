import psycopg2
from random import random
from time import time

NUMBER = 100
DATABASE = "property"
prices = [int(random() * 100000) for i in range(NUMBER)]
conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
cur = conn.cursor()
cities = ["Samarkand", "Bukhara", "Tashkent", "Andijan", "Nukus", "Fergana", "Navoi", "Termez", "Qarshi", "Kokand"]
types = ["House", "Apartment", "Dacha"]

def get_args(property):
        return property['id'], property['price'], property['isForSale'], \
            property['numberOfRooms'], property['location'], property['address'],\
            property['internalArea'], property['title'], property['description'], \
            property['publishDate'], property['authorID'] ,property['orienter'], property['propertyType'], property['isPublished']

def insert_dummy_property_from_dict(p, sql_statment):
    cur.execute(sql_statment, get_args(p))
    cur.execute("SELECT nextval('propertyimages_sequence')")
    _id=cur.fetchone()[0]
    conn.commit()
    insert_images(p['price'], _id)
    return _id

def delete_properties():
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("DELETE FROM Properties WHERE id NOT IN(SELECT PropertyID FROM PropertyImages)")    
    conn.commit()
    conn.close()

def insert_images(price, i):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT id FROM Properties WHERE price=%s"%price)
    _id  = cur.fetchone()[0]
    
    cur.execute("INSERT INTO PropertyImages(id, propertyID, filename) VALUES(%s , %s,'./assets/property_images/default.png')"%(i, _id))
    
    conn.commit()
    conn.close()

property = {"id" : 0, "price" : 391, "isForSale": True, "numberOfRooms": 4,
         "location" : "Samarkand", "address" : "sa","internalArea" : 200, 
         "title": "Some title", "description": "some description",
         "publishDate" : "2022-01-01", "authorID" : 1, "orienter" : "At the Center", "propertyType" : "House",
        "isPublished":True
         }


def main():
    start = time()
    for i in range(10, NUMBER):
        property['id'] = i
        property['location'] = cities[i%10]
        property['numberOfRooms'] = int(random() * 5) + 1
        property['propertyType'] = types[i%3]
        property['price'] = prices[i]
        l = str("%s, ")* 14
        sql_statment = """INSERT INTO Properties(
        ID, price, isForSale, numberOfRooms, location,
        address, internalArea, title, description, 
        publishDate, authorID, orienter, propertyType, isPublished) 
        VALUES(""" + l[:-2]  +  ')'
        _id  = insert_dummy_property_from_dict(property, sql_statment)
        print(_id)
    delete_properties()

    print(time() - start, " Seconds took to insert ", NUMBER , " Property" )
main()