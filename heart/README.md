# github.com/yindaheng98/gogistry/heart

This two package is the "heart" of higher-level protocol in gogistry (lower-level protocol should be implemented by yourself in the package `Protocol`). It contains the only way to run gogistry system (the only way to "beat" the "heart").

This package was designed for the controlling of the "heartbeat" sequence in gogistry. The "heartbeat" means a list of registrant information is keeping in each registry, and can be accessed from some interface; registrant should frequently send their information to registry ("update" the "connection"), or their previous information will be deleted from the registrant list ("disconnection"). This package is going to control the frequency of "heartbeat" sending from registrant, and decide when should a registry delete a registrant information from the registrant list.

A usage example is in [example/heart](../example/heart)