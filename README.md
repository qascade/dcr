# dcr
A library to orchestrate interoperable Data Cleaning Services using Intel SGX based Containers.

This Repo follows following [Contributing Guidelines](https://github.com/qascade/dcr/blob/main/CONTRIBUTING.md)

NOTE: 
1. As of now this library is only able to show a Proof of Concept for the architecture in the spec. 
2. Current PoC only shows an example of three collaborators, out of which two provide sources and one provides the transformation. 
3. Currently only Confidential GoApps are supported for transformations.
4. The Code is not production ready and does not partake any security measures other than access control and differential privacy. 
5. The library is still not tested on actual SGX backed machines and but the PoC can be tested on simulation mode. 

Links: 
1. [Spec Doc](https://www.notion.so/Clean-Room-Spec-Doc-f606d90163ff4ca9b14bae92c0db328d?d=78e16509ae124e7db6777a751a72cbb3#6e16fc663e0147f6b844030c4ce9fac8)
2. [Implementation Plan](https://www.notion.so/Implementation-Plan-e105e6e1a2d94d4ba6547cab5705ab20?pvs=4)

