#include "hashserver_table.hpp"
#include <iostream>

void hashserver_table::put(std::string key, std::string value)
{   
    std::lock_guard<std::mutex> lck(lock); 
    key_value_table[key] = value;
}

std::string hashserver_table::get(std::string key)
{   
    std::string value_found;
    auto search_key = key_value_table.find(key);
    if (search_key == key_value_table.end())
        return " ";
    else {
        value_found = search_key->second;
        std::cout << "value: " << value_found << std::endl;
        return value_found; 
    }
}

int main(int argc, char* argv[]) {
    std::cout << "Hello world" << std::endl;
    hashserver_table h;
    h.put("1","test");
    h.get("1");
    return 0;
}