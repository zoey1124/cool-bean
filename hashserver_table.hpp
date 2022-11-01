#include <map>
#include <mutex>

class hashserver_table{
public:
    void put(std::string key, std::string value);
    std::string get(std::string key);
private:
    std::map<std::string, std::string> key_value_table;
    std::mutex lock; 
};