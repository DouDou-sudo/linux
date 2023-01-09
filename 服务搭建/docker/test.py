import os
class CLanguage :
    # 下面定义了2个类变量
    name = "C语言中文网"
    add = "http://c.biancheng.net"
    def __init__(self,name,add):
        #下面定义 2 个实例变量
        self.name = name
        self.add = add
        print(name,"网址为：",add)
    # 下面定义了一个say实例方法
    def say(self, content):
        print(content)
# 将该CLanguage对象赋给clanguage变量
clanguage = CLanguage("C语言中文网","http://c.biancheng.net")

    #输出name和add实例变量的值
print(clanguage.name,clanguage.add)
    #修改实例变量的值
# clanguage.name="Python教程"
# clanguage.add="http://c.biancheng.net/python"
    #调用clanguage的say()方法
# clanguage.say("人生苦短，我用Python")
#     #再次输出name和add的值
# print(clanguage.name,clanguage.add)
# os.path.exists()
x=1
def a():
    # x +=1
    print(x)
a()