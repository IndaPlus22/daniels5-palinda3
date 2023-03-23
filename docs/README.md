If written answers are required, you can add them to this file. Just copy the
relevant questions from the root of the repo, preferably in
[Markdown](https://guides.github.com/features/mastering-markdown/) format :)

## Task 1
1. If i remove the go command from seek, then the program will run the seek function sequentially. This means that the person recieve a message will always be from the person to the left in the array. This is because the program will launch first "Anna" too seek someone to recieve after that it will always launch "Bob" and Bob will see one name in the channel and thus recieve the message from Anna and so forth.

2. If make the suggested changes the program will not compile but instead get a deadlock error. This is because the changing of "wg := new(sync.WaitGroup)" to "var wg sync.WaitGroup" will make the wg varaible contain the sync.WaitGroup struct itself, the old decleration created a pointer to the sync.waitGroup struct. Furthermore if we now make the last change and remove "*". That tells the functions that it wants a direct copy of the passing argument instead of a pointer. This makes it so that when we pass in wg from main function it will create a copy of wg to be used in the seek functions. Since the wg is a copy the changes made in the seek functions will not affect the wg in main, so it will forever wait for the seek functions to be done even when they are so.

3. If we remove the buffer of 1, then the behavour of how the sender blocks the program. A absent buffer results in the sender blocking until a recieve has recived the value. This causes a deadlock in this program since the last person will wait forever until a another seek function is avalible to recieve the value, but since no other seek function exept the last one is running it will never be recieved and thus deadlock. However if we have the buffer then the last person can safley send the value aslong as the capacity is not reached and continue, then when a recieving function AKA the main in the last case is ready to recieve it will do so.

4. If we remove the default from the main select. That could become a issue if we have a even array. In the case that all people get their message recieved and sent there will be no leftover message in the channel this will cause the select function in the main to forever wait for a message causing a deadlock. However with the default in it can automaticly exit the select if there is no avalible message to be recieved or taken from the buffer.

|Variant       | Runtime (ms) |
| ------------ | ------------:|
| singleworker |       709 ms    |
| mapreduce    |          310 ms |