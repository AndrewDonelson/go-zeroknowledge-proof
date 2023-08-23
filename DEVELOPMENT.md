Switching to symmetric cryptography can indeed simplify key management and reduce computational overhead. Let's look at how to adapt your idea using symmetric cryptographic techniques, and then consider ways to bolster data security:

### Symmetric Cryptographic Approach:

1. **Shared Secret**: Every authorized node in the network should share a secret key with Node B (the proving node). This secret key will be used for encryption and decryption. This means Node B will have a separate shared key for every other node in the network.

2. **Authorization**: 
   - Node A sends a request to Node B.
   - Node B checks if it has a shared key for Node A (i.e., Node A is authorized).

3. **Proof Generation with UUID**:
   - Node B combines the requested data with Node A's UUID.
   - Generates a hash of the combined data.
   - Encrypts the hash using the shared secret key with Node A.

4. **Proof Verification**:
   - Node A decrypts the received proof using the shared secret key.
   - Verifies the proof by hashing the original data and its UUID, and then comparing it to the decrypted proof.

### Enhancing Data Security:

1. **Use Secure Symmetric Algorithms**: Algorithms like AES (Advanced Encryption Standard) in GCM (Galois/Counter Mode) or CTR (Counter Mode) are considered secure and are widely adopted.

2. **Nonce/Initialization Vector (IV)**: Always use a new nonce or IV for each encryption operation. This ensures that even if the same data is encrypted multiple times, the ciphertext will be different. The nonce/IV can be sent alongside the ciphertext, as it's not a secret value but should be unique for each encryption operation with the same key.

3. **Key Rotation**: Periodically rotate the shared secret keys. While symmetric key management is simpler than asymmetric, key rotation is still a good practice to ensure long-term security.

4. **Authenticated Encryption**: Use authenticated encryption modes (like AES-GCM) that provide both confidentiality and integrity. This way, Node A can be sure that the proof received hasn't been tampered with.

5. **Rate Limiting & Monitoring**: Implement rate limiting on Node B to prevent brute-force or DDoS attacks. Monitor requests to identify and block any malicious activity.

6. **HMAC for Additional Integrity**: Along with encryption, you can use HMAC (Hash-Based Message Authentication Code) with a separate secret to ensure the integrity and authenticity of the message. This might be redundant if you're using authenticated encryption but can be an added layer of security in other cases.

7. **Salting the Data**: Before hashing the combination of data and UUID, introduce a salt. This ensures that even if two nodes request proof for the same data, the resulting hashes (and thus the encrypted proofs) will be different.

8. **Secure Key Distribution**: While symmetric key distribution is simpler than asymmetric, it's essential to ensure that shared secret keys are distributed securely, perhaps through a secure channel or an out-of-band mechanism.

Remember, the strength of a cryptographic system often lies not just in the algorithms used but also in its implementation and the protocols surrounding it. Testing and expert review are crucial before deploying such a system in a real-world scenario.

## Data  Redundancy

Broadcasting data to all authorized nodes and having them verify the data before accepting it provides a solid foundation for a redundant and consistent system. 

### 1. Data Addition and Broadcast:

When a node adds new data:

1. It computes a proof (hash) of the data.
2. It broadcasts the data along with its proof to all other authorized nodes. 

**Endpoint: `/broadcast`**

- **Method:** `POST`
- **Payload:** 
  - `data`: The new data to be shared.
  - `proof`: The cryptographic hash of the data (from the `/prove` endpoint).

### 2. Data Verification:

Upon receiving the broadcasted data, each node:

1. Computes the hash of the received data using the `/prove` endpoint.
2. Compares the computed hash with the received hash using the `/verify` endpoint.
3. If the hashes match, the data is accepted and added to the node's storage.

This can be an internal operation and might not need a dedicated endpoint unless you want external visibility or control over this process.

### 3. Data Acceptance:

If a node verifies and accepts the data:

1. It sends an acknowledgment back to the originating node.
2. Optionally, it can also broadcast its acceptance to other nodes, providing additional verification for the network.

**Endpoint: `/acknowledge`**

- **Method:** `POST`
- **Payload:** 
  - `dataId`: An identifier for the acknowledged data (could be a UUID or hash).
  - `nodeId`: Identifier of the acknowledging node.
  - `status`: Acceptance status (e.g., `accepted`, `rejected`).

### Suggested Improvements:

1. **Batching:** Instead of broadcasting individual pieces of data, you could batch multiple data entries together and broadcast them at once. This reduces network overhead.
 
2. **Conflict Resolution:** Have a strategy in place for handling situations where a node rejects the data. This could involve re-broadcasting, requesting data from multiple nodes to achieve consensus, or manual intervention.

3. **Secure Broadcast:** Ensure that the broadcast mechanism is secure. You might want to consider encrypting the data during transmission and using secure channels (like HTTPS or secure websockets) for communication.

4. **Node Trustworthiness:** Consider a mechanism to evaluate the reliability or trustworthiness of nodes. Nodes that consistently provide bad data or fail verifications could be flagged or removed from the network.

5. **Data Retention and Archiving:** As data grows, consider strategies for archiving old data, ensuring that all nodes don't need to store the entire dataset indefinitely.

Implementing these endpoints and strategies should provide a robust mechanism for data redundancy across nodes, ensuring data integrity and availability.