# dcr
A proof of concept framework to orchestrate Interoperable Differentially Private Data Clean Room Services on Intel SGX.

A Data Clean Room is a secure environment where organizations can collect data from multiple sources and combine it with their first-party data. Doing so allows marketers to leverage large, aggregated datasets of consumer behavior to provide insight into critical factors like performance, demographics, campaigns, etc.

*Data clean rooms allow companies to extract value from aggregate datasets sourced from multiple parties while prioritizing user privacy and maintaining strict security measures.*
## Contributing Guidelines
This Repo follows following [Contributing Guidelines](https://github.com/qascade/dcr/blob/main/CONTRIBUTING.md)

## NOTE: 
1. As of now this framework is only able to show a Proof of Concept for the architecture described in the spec. 
2. Current PoC only shows an example of three collaborators, out of which two provide sources and one provides the transformation. 
3. Currently only Confidential GoApps are supported for transformations.
4. The Code is not production ready and does not partake any security measures other than access control and differential privacy. 
5. The library is still not tested on actual SGX backed machines and but the PoC can be tested on simulation mode. 

To Run an example Data Clean Room Scenario.  
1. Make sure you have [Intel-SGX](https://github.com/intel/linux-sgx-driver) SDK installed. 
2. Make sure you have [Go](https://go.dev/) 1.20+ installed. 
3. Make sure you have [E-Go](https://github.com/edgelesssys/ego) Compiler installed on your PC's

You can see the sample collaboration package in samples/init_collaboration. You can also look at some sample packages under construction along with their graph images in samples/test_graph

## To Run the demonstration: 
1. Build the `dcr` binary.
```bash 
make build 
```
<img width="1033" alt="Screenshot 2023-05-11 at 4 09 17 AM" src="https://github.com/qascade/dcr/assets/92882879/30ae2382-6775-405d-9902-adafdb764251">

<img width="1119" alt="Screenshot 2023-05-11 at 4 10 33 AM" src="https://github.com/qascade/dcr/assets/92882879/8cd81896-cfd8-4724-85aa-afd7b2829c9b">

2. Run the demonstration 
```bash
./dcr run --pkgpath samples/init_collaboration
```

Links: 
1. [Spec Doc]( https://cliff-colt-e2a.notion.site/Solution-3-07d81059daab40cb84180336a33c3dd9)
2. [Research Doc](https://cliff-colt-e2a.notion.site/Clean-Room-Doc-f606d90163ff4ca9b14bae92c0db328d)

