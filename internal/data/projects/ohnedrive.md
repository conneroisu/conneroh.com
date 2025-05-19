---
id: ohnedrive
aliases:
  - OhneDrive
tags:
  - programming-language/c
  - programming-language/python
  - programming-language/cpp
created_at: 2025-03-27T14:13:11.000-06:00
description: OhneDrive is a open source self-driving mini-bus system that can be used to transport people and goods with an underlying purposal of data collection.
title: OhneDrive
updated_at: 2025-05-16T10:43:14.000-06:00
---

# OhneDrive

OhneDrive is a open source self-driving mini-bus system that can be used to transport people and goods with an underlying purposal of data collection.


Proposal presented at the Move the World Competition at the Innovation Center at Iowa State University, Ames, Iowa, March 24, 2022.

**Conner D. Ohnesorge[^1]**
Department of Electrical Engineering
Iowa State University
Ames, IA 50011
[connero@iastate.edu](mailto:connero@iastate.edu)
March 24, 2024

### ABSTRACT

CyRide is the public transportation system in Ames, Iowa. The system is a vital part of the community, providing transportation to students, faculty, and staff at Iowa State University, as well as to the public. The system is also a key part of the city’s infrastructure, providing a means of transportation for those who do not have access to a car. For insurance purposes, CyRide places GPS tracking devices, 360-degree cameras, accelerometers, and other sensors on all of its buses. However, the data collected from these sensors is not currently being used to its full potential as it is deleted after a certain period of time, generally 30 days. This proposal outlines a plan to collect the existing data on the bus system in Ames, Iowa, with the goal of automating the system in the future.

**Keywords:** CyRide · Ames, Iowa · public transportation · autonomous vehicles · data collection

---

## 1. Introduction

In the 21st century, the concept of an AI winter, a period of no advancements within the field of AI, is widely regarded as having passed, as referenced in academic literature (Toosi et al. \[2021]), with artificial intelligence now integrated into our daily lives and seemingly limitless new research findings. Among the forefronts of AI application is the domain of self-driving vehicles. This technology is evolving at a brisk pace, and it stands to reason that self-driving cars will soon become commonplace in our lives. Although the challenge of developing autonomous vehicles is formidable, leveraging data significantly facilitates this process. Without this data, the development of autonomous vehicles is currently deemed impossible. Data collection is a crucial step in the development of current autonomous vehicles.

> **Figure 1.** CyRide Bus During the Fall

The cruciality of data collection is evident in the fact that the most advanced autonomous vehicles currently being tested in the real world are restricted to certain areas (Sun et al. \[2020]). Some examples of companies that are currently testing autonomous vehicles are Tesla, Waymo, and Uber. Each of these companies has collected vast amounts of data to train their AI models. Additionally, all the companies’ vehicles are largely limited to certain areas, such as Waymo’s vehicles being limited to the Phoenix area (Sun et al. \[2020]). While Tesla’s Autopilot system is currently the most advanced, it is still not fully autonomous; they have become the most advanced by collecting the most data (Li et al. \[2024]).

Simply put, the more data that is available, the better the AI model can be trained in the specific environment. While potentially a splitting issue because of privacy, the data collected by CyRide can be used to train an AI model to drive fleets of the CyRide buses. The adoption of this proposed data collection effort would be a great benefit to the community, further establishing Iowa State University as a leading engineering school, and could provide more jobs to the community.

Without the adoption of this proposal, public transportation in Ames, Iowa, would not be able to take advantage of the benefits of autonomous vehicles, and the community will be left behind in the development of this technology.

---

## 2. Project Objectives and Goals

The primary objective of this proposal is to leverage the extensive data collected by the CyRide bus system in Ames, Iowa, to lay the groundwork for future autonomous vehicle integration within the city’s public transportation system. Specifically, this project aims to:

* Analyze the current data collection methodologies employed by CyRide, identifying key data points relevant to autonomous vehicle technology.
* Develop a comprehensive data storage solution that extends beyond the current 30-day retention period, ensuring valuable data is preserved for long-term analysis and research.
* Implement advanced data analysis techniques to understand traffic patterns, bus route efficiency, and potential challenges for autonomous vehicles in an urban environment.
* Explore the feasibility of introducing autonomous vehicle technology to Ames’ public transportation system, focusing on safety, efficiency, and scalability.
* Establish a collaborative framework for ongoing research between Iowa State University and CyRide, promoting innovation in public transportation and autonomous vehicle technology.
* Further establish Iowa State University as a leading engineering school.

---

## 3. Benefits

Accident prevention is a major benefit of the proposed partial adoption of autonomous vehicles for public transportation in Ames. According to the National Highway Traffic Safety Administration, 94% of accidents are caused by human error. Of all accidents, CyRide reports that 5 of 22 accidents occurred with stationary objects while the remaining occurred within their own facilities (Daily \[2024]).

### 3.1 Qualifications

As a student at Iowa State University, I’m well-equipped with the resources required to undertake this project. My access to CyRide buses and their collected data, alongside the comprehensive resources available at the university such as the library and faculty expertise, form the foundation of my project capabilities. My skillset, which encompasses programming and data analysis, is bolstered by my academic coursework, professional experiences, and personal ventures. A testament to my proficiency in handling AI-related data tasks is my contributions to and position (Chief AI Officer) in KreativeDocuvet Organization. In this startup, which I co-founded, I developed a custom AI transcription system that includes diarization and a Large Language Model, showcasing my ability to innovate in the field of artificial intelligence.

---

## 4. Methodology

This section outlines the proposed methodology for collecting, storing, and analyzing data from the CyRide bus system in Ames, Iowa. Our approach consists of three main stages: Data Collection, Data Storage, and Data Analysis.

### 4.1 Data Collection

This proposal outlines a plan to collect data from a variety of sensors already installed on CyRide buses, including GPS devices, 360-degree cameras, and accelerometers. Data collection will be continuous, providing a rich dataset that captures the dynamics of the Ames transit environment in real time. Key data points will include bus location, speed, acceleration, deceleration, passenger load, and steering wheel position, among others.

Recently, CyRide even upgraded the cameras on their buses and facilities to higher-quality cameras (Tribune \[2024]). The plan includes:

* Establishing a data pipeline from CyRide’s existing infrastructure to our secure, high-capacity storage solutions.
* Ensuring data integrity and confidentiality through state-of-the-art encryption and access control mechanisms.
* Implementing data processing workflows that prepare raw sensor data for analysis, including cleaning, normalization, and preliminary analysis.

### 4.2 Data Storage

All collected data will be stored securely in a cloud-based storage solution, ensuring easy access for analysis while maintaining strict data privacy and security protocols. We propose using the scalable storage Amazon S3, allowing us to adjust storage resources as the dataset grows. The adoption of Amazon S3 is proposed for several reasons:

* Complies with industry standards for data privacy and security.
* Ensures scalability and reliability as the dataset grows.
* Balances long-term retention with privacy and compliance requirements.
* Leverages the vast ecosystem of tools available for S3-based workflows.

Amazon S3 provides 99.999999999% durability and 99.99% availability (Alliance \[2010]), and despite past misconfigurations, remains a secure service (Thompson \[2019]).

### 4.3 Data Upload

The data collected from the CyRide buses will be uploaded to Amazon S3 via a secure pipeline, ensuring encryption in transit and at rest. Uploads will occur in real time, providing immediate access to the latest bus location, speed, and sensor data.

### 4.4 Data Analysis

Analysis will combine traditional statistical methods with advanced machine learning algorithms. Camera feeds will also be used to predict steering wheel orientation. Key focus areas include route optimization, traffic pattern analysis, and safety enhancement strategies.

---

## 5. Challenges and Considerations

1. **Privacy Issues:** Protecting the privacy of individuals captured in sensor data is paramount and requires strict access controls and anonymization strategies.
2. **Technological and Infrastructure Investment:** Significant funding and infrastructure upgrades will be needed; potential sources include grants, partnerships, and private investment.
3. **Cultural and Social Acceptance:** Community engagement and transparent communication are essential to address concerns about job displacement, safety, and trust in autonomous systems.

---

## 6. Expected Outcomes

* A comprehensive, long-term dataset to advance autonomous-vehicle research.
* Immediate insights for improving current CyRide operations (efficiency and safety).
* Educational opportunities for ISU students in data science and AI.
* A roadmap toward partial integration of autonomous buses on select routes, reducing wait times, lowering costs, and decreasing human-error accidents.

---

## 7. Timeline and Budget

**Timeline (2 months total):**

1. Data Collection System Creation (API → S3): 1 month
2. Data Storage Implementation (S3): 0.5 month
3. Data Analysis Setup & Validation: 0.5 month

Initial budget covers expanded storage, software licenses, and personnel.

---

## 8. Estimated Expenses

The estimated monthly expenses amount to \$350.

**Table 1: Estimated Monthly Expenses for Data Collection Project**

| Expense Category      | Estimated Cost (USD) |
| --------------------- | -------------------- |
| Amazon S3 Storage     | \$100                |
| API Hosting           | \$100                |
| Data Collection Tools | \$150                |

**Table 2: Infrequent, Low-Cost S3 Storage Options. *Proposed option underlined.***

| Storage Option                                    | Cost (per GB / Month) |
| ------------------------------------------------- | --------------------- |
| <u>S3 Intelligent – Infrequent Access Tier</u>    | \$0.0125              |
| <u>S3 Intelligent – Infrequent Access Tier</u>    | \$0.0125              |
| S3 Intelligent – Archive Instant Access Tier      | \$0.004               |
| S3 Intelligent – Tiering Archive Access Tier      | \$0.0036              |
| S3 Intelligent – Tiering Deep Archive Access Tier | \$0.00099             |
| S3 Standard – Infrequent Access                   | \$0.0125              |
| S3 Glacier Instant Retrieval                      | \$0.004               |
| S3 Glacier Flexible Retrieval                     | \$0.0036              |
| S3 Glacier Deep Archive                           | \$0.00099             |
| S3 One Zone – Infrequent Access                   | \$0.01                |

**Table 3: Pricing for Various AWS S3 Services. *Asterisk denotes proposed service.***

| Service                  | PUT, COPY, POST, LIST (\$/1,000 req.) | GET, SELECT, etc. (\$/1,000 req.) | Lifecycle Transition (\$/1,000 req.) |
| ------------------------ | ------------------------------------- | --------------------------------- | ------------------------------------ |
| S3 Standard              | \$0.005                               | \$0.0004                          | n/a                                  |
| S3 Intelligent-Tiering\* | \$0.005                               | \$0.0004                          | \$0.01                               |

---

## 9. Conclusion

This proposal outlines a visionary plan to harness CyRide’s unused sensor data to improve current operations and pave the way for autonomous buses in Ames. Aligning with Iowa State University’s innovation mission, this project can serve as a national model for data-driven public transportation.

---

## References

1. Toosi, A., Bottino, A., Saboury, B., Siegel, E.L., & Rahmim, A. (2021). *A brief history of AI: how to prevent another winter (a critical review)*. CoRR, abs/2109.01517. [https://arxiv.org/abs/2109.01517](https://arxiv.org/abs/2109.01517)
2. Sun, P., Kretzschmar, H., Dotiwalla, X., Chouard, A., Patnaik, V., Tsui, P., Guo, J., Zhou, Y., Chai, Y., Caine, B., Vasudevan, V., Han, W., Ngiam, J., Zhao, H., Timofeev, A., Ettinger, S., Krivokon, M., Gao, A., Joshi, A., Cheng, S., Zhang, Y., Shlens, J., Chen, Z., & Anguelov, D. (2020). *Scalability in perception for autonomous driving: Waymo open dataset*.
3. Li, L., Shao, W., Dong, W., Tian, Y., Zhang, Q., Yang, K., & Zhang, W. (2024). *Data-centric evolution in autonomous driving: A comprehensive survey of big data system, data mining, and closed-loop technologies*.
4. Iowa State Daily. (2024). *Breaking down CyRide’s recent accident totals*. [https://iowastatedaily.com/272090/uncategorized/breaking-down-cyrides-recent-accident-totals/](https://iowastatedaily.com/272090/uncategorized/breaking-down-cyrides-recent-accident-totals/)
5. Ames Tribune. (2024). *Here’s what you may have missed at the latest Ames city council meeting*. [https://www.amestrib.com/story/news/local/2024/01/24/heres-what-you-may-have-missed-at-the-latest-ames-city-council-meeting/72324987007/](https://www.amestrib.com/story/news/local/2024/01/24/heres-what-you-may-have-missed-at-the-latest-ames-city-council-meeting/72324987007/)
6. Cloud Security Alliance. (2010). *Amazon AWS – 11 9’s of reliability?*. [https://cloudsecurityalliance.org/blog/2010/05/24/amazon-aws-11-9s-of-reliability](https://cloudsecurityalliance.org/blog/2010/05/24/amazon-aws-11-9s-of-reliability)
7. Thompson, P.A. (2019). *Capital One data breach*. [https://www.justice.gov/usao-wdwa/page/file/1194001/dl?inline](https://www.justice.gov/usao-wdwa/page/file/1194001/dl?inline)
8. Amazon Web Services. (2022). *Amazon S3 Simple Storage Service Pricing*. [https://aws.amazon.com/s3/pricing/](https://aws.amazon.com/s3/pricing/)

[^1]: Data Collection Proposal for CyRide in Ames, Iowa.
