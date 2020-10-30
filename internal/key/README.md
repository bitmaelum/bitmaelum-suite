We allow different kind of keys to be stored. Examples are API keys, and Auth keys. In practice, these keys are very
similar and creating their own storage systems for them would result in lots of duplicate code (which our CI system 
does not approve of).

So instead, we have a generic "key" system that allows us to store multiple types of keys, each with their own 
properties. This however, results in a more complex setup, as Go isn't able to do OO.


    +-------------------+           +------------------+          +-----------+
    | ApiKeyRepository  |           | (storageBackend) |          | boltRepo  |
    |   (ApiKeyRepo)    | ------+-> |                  | -+-----> |           |
    +-------------------+       |   +------------------+  |       +-----------+
                                |                         |
                                |                         |
                                |                         |
    +-------------------+       |                         |       +-----------+
    | AuthKeyRepository |       |                         |       | mockRepo  |
    |   (AuthKeyRepo)   | ------+                         +-----> |           |
    +-------------------+                                 |       +-----------+
                                                          |       +-----------+
                                                          +-----> | redisRepo |
                                                                  +-----------+
                                                                 
Both the ApiKeyRepository and AuthKeyRepository deal with the corresponding ApiKeyType and AuthKeyType structures. 
They both are a small wrapper around a storage repository that deals with interfaces{}. So, we convert these 
interfaces to actual correct types in the api and key repository:

``` 
func (a AuthKeyRepository) Fetch(ID string) (AuthKeyType, error) {
    v := &AuthkeyType{}
	err := a.repo.Fetch(ID, v)
	return v, err
}
```

The storage repository will simply fetch data from the storage, and unmarshals this trough the given `v` value. 
Sometimes we need to use GenericKey functionality (`GetID()`, `GetHashAddress()`), which is where we typecast this
at certain places in the storage repositories.

There is a bit of reflection going on when dealing with lists. This is because we need to allocate items, but we don't 
know which items. We deal with this in the `Fetch` methods by passing the actual storage value, but we cannot do this
in the `FetchByHash` method. Instead, we pass again, a single structure (like in `Fetch()`), and based on that value,
we create with reflection a new variable. This variable gets marshalled through JSON and stored in the list.

Finally, the ApiKey and AuthKey repositories will retrieve these lists with undefined interfaces, and cast them to
either `[]ApiKeyType` or `[]AuthKeyType`.

With the amount of effort it took to get this up and running decently, DRY isn't always holy. 
