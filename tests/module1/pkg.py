import module1.submodule1.pkg
from module1.submodule1 import pkg as sub

def hello_world():
    print("this is module1.pkg")

sub.hello_world()
