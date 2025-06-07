# Limitations

- Currently, files are read fully into memory using `os.ReadFile`. This approach is not suitable for very large files, as it can cause high memory usage or crashes.

# Improvements
- When directories are sent, we auto zip and then transfer
- implement chunked reading/streaming to handle large files efficiently. 
- Sender is able to set a key with send flag which will be used to encrypt and decrypt
- Make these functions accessable through REST APIs
- Make mobile first web app to make the product more well rounded