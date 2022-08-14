import re
from create_database import create_tables
import requests
import unittest
import json
import psycopg2

DATABASE = 'property'
URL = "http://127.0.0.1:8000"


class TestA(unittest.TestCase):
    @classmethod
    def tearDownClass(cls):
        pass
    @classmethod
    def setUpClass(self):
        create_tables()
        self.user = {"id": 0, "email": "mehemtmuratkafadaroglu@gmail.com",
         "companyName": "Some", "isAgent": False, 
         "phoneNumber": 7388807595, 'password' : "123456789", 'isMailVerified' : False}
        print("\nTesting User")        
    def is_user_verified(self, _id):
        res = get_user(" WHERE ID = %d"%_id)[0][6]
        return res

    def test_insert(self):
        requests.post(URL+"/public/user/", json=self.user)
        res = get_user(" WHERE companyName = 'Some'")
        res = res[0]
        _id = res[0]
        verify_user(_id, get_temp_key(_id)) #verify user
        mail = res[1]
        company_name = res[2]
        is_agent = res[3]
        number= res[4]
        self.assertEqual(mail, self.user['email'])
        self.assertEqual(company_name, self.user['companyName'])
        self.assertEqual(is_agent, self.user['isAgent'])
        self.assertEqual(number, self.user['phoneNumber'])
        self.assertTrue(self.is_user_verified(_id))

    def test_select(self):
        pass
    def test_delete(self):
        pass
    def test_edit(self):
        pass

class TestB(unittest.TestCase):
    @classmethod
    def tearDownClass(self):
        pass    
    @classmethod
    def setUpClass(self):
        self.user = {"id": 0, "email": "mehemtmuratkafadaroglu@gmail.com",
         "companyName": "Some", "isAgent": False, 
         "phoneNumber": 7388807595, 'password' : "123456789", 'isMailVerified' : False}
        self.user_mail = self.user['email']
        self.user_password = self.user['password']
        self.token = login(self.user_mail, self.user_password)
        self.userID = get_user_id("  WHERE email='%s'"%self.user_mail)
        self.property = {"id" : 0, "price" : 391, "IsForSale": True, "numberOfRooms": 4,
         "location" : "Samarkand", "address" : "sa","internalArea" : 200, 
         "title": "Some title", "description": "some description",
         "publishDate" : "2022-01-01", "authorID" : self.userID, "orienter" : "At the Center", "propertyType" : "House"
         }
        print("\nTesting Property")

    def test_insert(self):
        requests.post(URL+"/private/add/properties/", 
        json=self.property, headers={"Authorization": self.token})
        conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
        cur = conn.cursor()
        cur.execute("SELECT * FROM Properties WHERE authorID = %d"%self.userID)
        res = list(cur.fetchall()[0])
        price = res[1]
        is_for_sale = res[2]
        number_of_rooms = res[3]
        location = res[4]
        address = res[5]
        internal_area = res[6]
        title =res[7]
        description = res[8]
        self.assertEqual(self.property['price'], price)
        self.assertEqual(self.property['IsForSale'], is_for_sale)
        self.assertEqual(self.property['numberOfRooms'], number_of_rooms)
        self.assertEqual(self.property['location'], location)
        self.assertEqual(self.property['address'], address)
        self.assertEqual(self.property['internalArea'], internal_area)
        self.assertEqual(self.property['title'], title)
        self.assertEqual(self.property['description'], description)
        conn.commit()
        conn.close()

    def test_select(self):
        url = "/public/properties/Samarkand/1000000/0/1000/0/1000/0"
        vals =get_as_json(url)
        properties = get_propeties(" WHERE location='Samarkand'")
        for i, property in enumerate(properties):
            _id = property[0]
            price = property[1]
            is_for_sale = property[2]
            number_of_rooms = property[3]
            location = property[4]
            address = property[5]
            internal_area = property[6]
            title =property[7]
            description = property[8]
            self.assertEqual(vals[i]['id'], _id)
            self.assertEqual(vals[i]['price'], price)
            self.assertEqual(vals[i]['isForSale'], is_for_sale)
            self.assertEqual(vals[i]['numberOfRooms'], number_of_rooms)
            self.assertEqual(vals[i]['location'], location)
            self.assertEqual(vals[i]['address'], address)
            self.assertEqual(vals[i]['internalArea'], internal_area)
            self.assertEqual(vals[i]['title'], title)
            self.assertEqual(vals[i]['description'], description)
    def test_delete(self):
        pass
    def test_edit(self):
        pass
    def test_save(self):
        conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
        cur = conn.cursor()
        cur.execute("SELECT id FROM Properties WHERE authorID = %d"%self.userID)
        temp = cur.fetchall()        
        property_id = temp[0][0]
        requests.get(URL + "/private/save/%s/%s"%(self.userID, property_id), headers={"Authorization": self.token})

def get_as_json(url):
    value = requests.get(URL + url).content
    return json.loads(value)

def get_user(where):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT * FROM Users" + where)
    res =  cur.fetchall()
    conn.commit()
    conn.close()
    return res

def get_user_id(where):
    return get_user(where)[0][0]
            
def get_propeties(where):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT * FROM Properties" + where)
    res =  cur.fetchall()
    conn.commit()
    conn.close()
    return res

def insert_dummy_user(email, password):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("""INSERT INTO Users(email,companyname,isagent,phonenumber,password,ismailverified)
                VALUES(%s, 'verified_user', false, 123123, %s, true);
    """, [email, password])
    conn.commit()
    cur.execute("SELECT ID FROM Users WHERE phonenumber=123123")
    res = cur.fetchone()[0]
    conn.close()
    return res

def get_temp_key(user_id):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT key FROM temporarykeys WHERE userID=%d"%user_id)
    res = cur.fetchone()
    conn.close()
    return res[0]

def verify_user(user_id, key):
    url ="/public/verify/%s/%s"%(user_id, key)
    requests.get(URL + url)

def get_token_from_db(user_id):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT token FROM Tokens WHERE userID=%d"%user_id)
    res = cur.fetchone()
    conn.close()
    return res[0]
    
def login(email, password):
    res = requests.post(URL + "/public/login/",json={"email": email,
            "password" : password
        })
    return json.loads(res.content)['token']

if __name__ == "__main__":
    token = login("mehemtmuratkafadaroglu@gmail.com","123456789")
    val = requests.get(URL+"/private/logout/", headers={"Authorization": token})
