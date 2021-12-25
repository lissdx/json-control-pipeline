
### About
This service implements JSON data pipeline process
The main idea is opportunity to create configurable data pipeline.

### Conf. file
Conf. file includes process name and conf. params:  
```
{
"steps":[  
   {
      "processor":"AddField",
      "configuration":{
         "fieldName":"firstName",
         "fieldValue":"George"
      }
   },
   {
      "processor":"AddField",
      "configuration":{
         "fieldName":"unexpected",
         "fieldValue":"unknown"
      }
   }]
}
```  
### How to add new Processor
1. Implement and add Process method to [Process](internal/pkg/process/json/process.go)
2. Update the [ProcessFactory](internal/pkg/process/json/process.go)

### Example
The example implementation of this data pipeline you can find in  
[main](cmd/main.go)

### Example Run result:
![MainResult](doc/Screen%20Shot%202021-12-25%20at%209.27.41.png)