from cgi import print_form
from create_database import create_tables
import requests
import unittest
import json
import psycopg2
import base64
import copy

DATABASE = 'property'
URL = "http://127.0.0.1:8000"


class Test(unittest.TestCase):
    @classmethod
    def setUpClass(self):
        create_tables()
        self.user = {"id": 0, "email": "mehemtmuratkafadaroglu@gmail.com",
         "companyName": "Some", "isAgent": False, 
         "phoneNumber": 7388807595, 'password' : "123456789", 'isMailVerified' : False}
        requests.post(URL+"/public/user/", json=self.user)
        self.res = get_user(" WHERE companyName = 'Some'")
        self.user_mail = self.user['email']
        self.user['id'] =  get_user_id("  WHERE email='%s'"%self.user_mail)
        self.userID =self.user['id']
        self.user_password = self.user['password']
        self.token = login(self.user_mail, self.user_password)
        self.property = {"id" : 0, "price" : 391, "isForSale": True, "numberOfRooms": 4,
         "location" : "Samarkand", "address" : "sa","internalArea" : 200, 
         "title": "Some title", "description": "some description",
         "publishDate" : "2022-01-01", "authorID" : self.userID, "orienter" : "At the Center", "propertyType" : "House",
         "images" : get_images(), "isPublished":True
         }
        self.properties = []
        print("\nTesting Property")

    def test_auser_insert(self):
        res = self.res[0]
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
        self.assertEqual(self.property['isForSale'], is_for_sale)
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
        properties = get_properties(" WHERE location='Samarkand'")
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
            self.assertEqual(vals[i]['property']['id'], _id)
            self.assertEqual(vals[i]['property']['price'], price)
            self.assertEqual(vals[i]['property']['isForSale'], is_for_sale)
            self.assertEqual(vals[i]['property']['numberOfRooms'], number_of_rooms)
            self.assertEqual(vals[i]['property']['location'], location)
            self.assertEqual(vals[i]['property']['address'], address)
            self.assertEqual(vals[i]['property']['internalArea'], internal_area)
            self.assertEqual(vals[i]['property']['title'], title)
            self.assertEqual(vals[i]['property']['description'], description)
    def test_delete(self):
        property = copy.deepcopy(self.property)
        _id = get_property_id()
        property['id'] = _id
        property['price'] = 123123
        insert_dummy_property_from_dict(property)
        requests.post(URL +"/private/delete/properties/", json=property,  headers={"Authorization": self.token})
        p = get_properties(" WHERE ID=%s"%_id)
        self.assertEqual(len(p), 0)
    def test_edit(self):
        insert_dummy_property(self.userID)
        property = copy.deepcopy(self.property)
        property['id'] = 2
        requests.post(URL+"/private/edit/properties/", 
        json=property, headers={"Authorization": self.token})
        res = get_properties(" WHERE ID=%s"%2)[0]
        location = res[4]
        self.assertEqual(location, property['location'])
        delete_dummy_property(self.userID)

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
        self.assertTrue(is_property_saved(self.userID, property_id))

    
    def test_select_users_properties(self):
        insert_dummy_property(self.userID)
        properties = get_properties(" WHERE authorID=%s"%self.userID)
        val = requests.get(URL + "/private/properties/%s"%(self.userID), headers={"Authorization": self.token})
        val  = json.loads(val.content)
        values = [i['property'] for i in val]
        for i, value in enumerate(values):
            property = properties[i]
            _id = property[0]
            price = property[1]
            is_for_sale = property[2]
            number_of_rooms = property[3]
            location = property[4]
            address = property[5]
            internal_area = property[6]
            title =property[7]
            description = property[8]
            self.assertEqual(value['id'], _id)
            self.assertEqual(value['price'], price)
            self.assertEqual(value['isForSale'], is_for_sale)
            self.assertEqual(value['numberOfRooms'], number_of_rooms)
            self.assertEqual(value['location'], location)
            self.assertEqual(value['address'], address)
            self.assertEqual(value['internalArea'], internal_area)
            self.assertEqual(value['title'], title)
            self.assertEqual(value['description'], description)
    def test_profile(self):
        res  = requests.get(URL +"/private/profile/",  headers={"Authorization": self.token})
        profile = res.content
        profile = dict(json.loads(profile))
        self.assertEqual(profile['email'], self.user['email'])
        self.assertEqual(profile['isAgent'], self.user['isAgent'])
        self.assertEqual(profile['phoneNumber'], self.user['phoneNumber'])
    def test_zedit_profile(self):
        user = copy.deepcopy(self.user)
        user['companyName'] = "GG"
        requests.post(URL +"/private/edit/profile/", json=user ,headers={"Authorization": self.token})
        res = get_user(" WHERE email ='%s'"%user['email'])
        company_name =res[0][2]
        self.assertEqual(company_name, "GG")

def is_property_saved(user_id, property_id):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT propertyID FROM Saved WHERE userID = %s"%user_id)
    res =  cur.fetchall()
    res = [i[0] for i in res]
    conn.commit()
    conn.close()
    return property_id in res

def insert_dummy_property_from_dict(property):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    def get_args():
        return property['id'], property['price'], property['isForSale'], \
            property['numberOfRooms'], property['location'], property['address'],\
            property['internalArea'], property['title'], property['description'], \
            property['publishDate'], property['authorID'] ,property['orienter'], property['propertyType'], property['isPublished']

    l = str("%s, ")* 14
    sql_statment = """INSERT INTO Properties(
    ID, price, isForSale, numberOfRooms, location,
    address, internalArea, title, description, 
    publishDate, authorID, orienter, propertyType, isPublished) 
    VALUES(""" + l[:-2]  +  ')'
    
    cur.execute(sql_statment, get_args())
    conn.commit()
    conn.close()

def insert_dummy_property(userID):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    def get_args():
        return get_property_id(), 100, True, 3, "Tashkent", "Some address", 100, "test Title", \
            "test description", "2022-01-01", userID, "test orienter", "House", True
    l = str("%s, ")* 14
    sql_statment = """INSERT INTO Properties(
    ID, price, isForSale, numberOfRooms, location,
    address, internalArea, title, description, 
    publishDate, authorID, orienter, propertyType, isPublished) 
    VALUES(""" + l[:-2]  +  ')'
    
    cur.execute(sql_statment, get_args())
    cur.execute(sql_statment, get_args())
    cur.execute(sql_statment, get_args())

    conn.commit()
    conn.close()

def delete_dummy_property(author_id):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("DELETE FROM Properties WHERE authorID=%s"%author_id)
    conn.commit()
    cur.close()

def get_property_id():
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT nextval('propertiesid_sequence')")
    res = cur.fetchone()[0]
    conn.close()
    return res
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
            
def get_properties(where, option="*"):
    conn = psycopg2.connect(host = "localhost", 
                            dbname=DATABASE, 
                            user="postgres",
                            password="630991")
    cur = conn.cursor()
    cur.execute("SELECT "  + option +" FROM Properties" + where)
    res =  cur.fetchall()
    conn.commit()
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

def get_images():
    with open("../assets/picture.jpg", 'rb') as f:
        s = base64.b64encode(f.read())
        val  = s.decode('utf-8')
    return [val]

if __name__ == "__main__":
    pass