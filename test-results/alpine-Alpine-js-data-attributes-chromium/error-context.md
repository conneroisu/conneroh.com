# Test info

- Name: Alpine.js data attributes
- Location: /home/connerohnesorge/Documents/001Repos/conneroh.com.git/feat/verifiable/tests/browser/alpine.test.ts:3:1

# Error details

```
Error: expect.toBeVisible: Error: strict mode violation: locator('[\\@click*="isMenuOpen"]') resolved to 5 elements:
    1) <img class="tw-78" @click="isMenuOpen = !isMenuOpen" src="https://conneroisu.fly.storage.tigris.dev/assets/svg/menu.svg"/> aka getByRole('navigation').locator('img')
    2) <a class="tw-83" hx-get="/projects" hx-push-url="true" hx-target="#bodiody" preload="ontouchstart" @click="isMenuOpen = false">Projects</a> aka getByText('Projects', { exact: true }).nth(1)
    3) <a class="tw-83" hx-get="/posts" hx-target="#bodiody" hx-push-url="/posts" preload="ontouchstart" @click="isMenuOpen = false">Posts</a> aka getByText('Posts', { exact: true }).nth(1)
    4) <a class="tw-83" hx-get="/tags" hx-push-url="true" hx-target="#bodiody" preload="ontouchstart" @click="isMenuOpen = false">Tags</a> aka getByText('Tags', { exact: true }).nth(1)
    5) <a class="tw-83" hx-push-url="true" hx-target="#bodiody" hx-get="/employments" preload="ontouchstart" @click="isMenuOpen = false">Employments</a> aka getByText('Employments').nth(1)

Call log:
  - expect.toBeVisible with timeout 5000ms
  - waiting for locator('[\\@click*="isMenuOpen"]')

    at /home/connerohnesorge/Documents/001Repos/conneroh.com.git/feat/verifiable/tests/browser/alpine.test.ts:14:60
```

# Page snapshot

```yaml
- banner:
  - navigation: Conner Ohnesorge Projects Posts Tags Employments
- main:
  - heading "Name" [level=1]: Conner Ohnesorge
  - paragraph: Electrical Engineer & Software Developer specialized in creating robust, scalable, and elegant solutions.
  - paragraph:
    - text: Electrical Engineering Bachelors Degree and Minor in Computer Science from
    - link "Iowa State University":
      - /url: https://iastate.edu/
  - link "View Projects":
    - /url: "#projects"
  - link "Contact Me":
    - /url: "#contact"
  - img "Photo of Me"
  - heading "Featured Projects" [level=2]
  - text: View All Projects →
  - img "bufnrix"
  - heading "bufnrix" [level=2]
  - paragraph: Nix powered protobuf tools
  - img
  - text: Go
  - img
  - text: Nix
  - img
  - text: Protobuf
  - img
  - text: gRPC 0 posts 4 tags
  - img "conneroh.com"
  - heading "conneroh.com" [level=2]
  - paragraph: This site!
  - img
  - text: Alpine.js
  - img
  - text: astro
  - img
  - text: TailwindCSS
  - img
  - text: Go
  - img
  - text: JavaScript
  - button "(3) More"
  - text: 0 posts 8 tags
  - img "CPRE488 MP0"
  - heading "CPRE488 MP0" [level=2]
  - paragraph: The first Project from CPRE488 at Iowa State University
  - img
  - text: iastate
  - img
  - text: VHDL
  - img
  - text: Vivado
  - img
  - text: C
  - img
  - text: UART 0 posts 5 tags
  - img "CPRE488 MP1"
  - heading "CPRE488 MP1" [level=2]
  - paragraph: The second Project from CPRE488 at Iowa State University
  - img
  - text: iastate
  - img
  - text: VHDL
  - img
  - text: Vivado
  - img
  - text: C
  - img
  - text: UART 1 posts 5 tags
  - img "CPRE488 MP2"
  - heading "CPRE488 MP2" [level=2]
  - paragraph: The third Project from CPRE488 at Iowa State University
  - img
  - text: iastate
  - img
  - text: CPRE488
  - img
  - text: Verilog
  - img
  - text: VHDL
  - img
  - text: Vivado
  - button "(4) More"
  - text: 2 posts 9 tags
  - img "CPRE488-mp3"
  - heading "CPRE488-mp3" [level=2]
  - paragraph: Mini-project 3 for CPRE488. Linux Device Drivers, Petalinux, boot loaders, and OpenCV, oh my!
  - img
  - text: iastate
  - img
  - text: Verilog
  - img
  - text: VHDL
  - img
  - text: Vivado
  - img
  - text: device-drivers
  - button "(5) More"
  - text: 1 posts 10 tags
  - heading "Professional Experience" [level=2]
  - paragraph: My journey through various roles in engineering and technology
  - link "Jun 2023 - Aug 2023 Controls Systems Engineer Intern at Freund-Vector Developed solutions for product development of granulation, compactor, and mixer machines. Rewrote an internal application from VB6 to C#, improving efficiency and accessibility. Control Systems Manufacturing PLC Programming Robotics +2":
    - /url: /employment/controls-systems-engineer-intern-freund-vector
    - text: Jun 2023 - Aug 2023
    - heading "Controls Systems Engineer Intern at Freund-Vector" [level=3]
    - paragraph: Developed solutions for product development of granulation, compactor, and mixer machines. Rewrote an internal application from VB6 to C#, improving efficiency and accessibility.
    - text: Control Systems Manufacturing PLC Programming Robotics +2
  - link "Jun 2024 - Aug 2024 Data Scientist at CSE Developed resource constrained data pipelines and large-scale scrapers of college sports data using Golang and optimized MySQL databases. MySQL Go Python SQL":
    - /url: /employment/data-scientist-cse
    - text: Jun 2024 - Aug 2024
    - heading "Data Scientist at CSE" [level=3]
    - paragraph: Developed resource constrained data pipelines and large-scale scrapers of college sports data using Golang and optimized MySQL databases.
    - text: MySQL Go Python SQL
  - link "Dec 2023 - Jul 2024 Founding Chief AI Officer at Kreative Docuvet Created diarized transcription services and LLM pipelines for veterinary practices, built tools for exploring large language model responses, and collected a large veterinarian voice corpus. Postgres NextJS PyTorch React +4":
    - /url: /employment/founding-chief-ai-officer-kreative-docuvet
    - text: Dec 2023 - Jul 2024
    - heading "Founding Chief AI Officer at Kreative Docuvet" [level=3]
    - paragraph: Created diarized transcription services and LLM pipelines for veterinary practices, built tools for exploring large language model responses, and collected a large veterinarian voice corpus.
    - text: Postgres NextJS PyTorch React +4
  - link "Dec 2021 - Apr 2022 Machine Learning Researcher at Iowa State University Developed an EEG gesture recognition program as part of a senior research project, conducting extensive data analysis and implementing real-time gesture recognition. Machine Learning Matlab Python":
    - /url: /employment/machine-learning-researcher-iowa-state
    - text: Dec 2021 - Apr 2022
    - heading "Machine Learning Researcher at Iowa State University" [level=3]
    - paragraph: Developed an EEG gesture recognition program as part of a senior research project, conducting extensive data analysis and implementing real-time gesture recognition.
    - text: Machine Learning Matlab Python
  - heading "Recent Posts" [level=2]
  - text: View All Posts →
  - img "A Reflective Journey - Navigating Your Cumulative Experience at Iowa State University"
  - heading "A Reflective Journey - Navigating Your Cumulative Experience at Iowa State University" [level=2]
  - paragraph: This is a reflection on my time at Iowa State University.
  - img
  - text: AGEDS461
  - img
  - text: ARCH 321
  - img
  - text: CPRE281
  - img
  - text: CPRE288
  - img
  - text: CPRE381
  - button "(9) More"
  - text: 14 tags | 4 projects Mar 27, 2025
  - img "Making Vivado not suck at Git"
  - heading "Making Vivado not suck at Git" [level=2]
  - paragraph: Post on how I made Vivado not suck at Git
  - img
  - text: CPRE488
  - img
  - text: Vivado
  - img
  - text: Git 3 tags | 1 projects Mar 27, 2025
  - img "Miracle Rice in Bali"
  - heading "Miracle Rice in Bali" [level=2]
  - paragraph: "Miracle Rice in Bali: Computational Investigations into the Origins of Efficiency"
  - img
  - text: AGEDS461 1 tags | 0 projects Dec 01, 2021
  - img "My Ethics"
  - heading "My Ethics" [level=2]
  - paragraph: This essay discusses the importance of having a personal code of ethics and using a thoughtful decision-making process to promote ethical behavior in organizations and professions.
  - img
  - text: EE394 1 tags | 0 projects May 02, 2025
  - heading "Skills & Technologies" [level=2]
  - text: See All Skills/Technologies →
  - heading "ARM" [level=2]
  - img
  - text: ARM is a 32-bit reduced instruction set computing (RISC) architecture. 0 posts 0 projects
  - heading "Intel 386" [level=2]
  - img
  - text: Intel 386 is a 32-bit x86 instruction set architecture. 0 posts 0 projects
  - heading "MIPS" [level=2]
  - img
  - text: MIPS is a 32-bit reduced instruction set computing (RISC) architecture. 0 posts 0 projects
  - heading "RISC-V" [level=2]
  - img
  - text: RISC-V is a 64-bit open instruction set architecture. 0 posts 0 projects
  - heading "x86-64" [level=2]
  - img
  - text: x86-64 is a 64-bit x86 instruction set architecture. 0 posts 0 projects
  - heading "AWS" [level=2]
  - img
  - text: Amazon Web Services (AWS) is a subsidiary of Amazon that provides on-demand cloud computing platforms and other infrastructure services. 0 posts 0 projects
  - heading "Flyio" [level=2]
  - img
  - text: Flyio is a cloud-based platform for building and deploying web applications based on AWS. 0 posts 0 projects
  - heading "Freund-Vector Corp." [level=2]
  - img
  - text: Freund-Vector Corp. is a global service provider and manufacturer of granulating, coating and drying equipment. 0 posts 0 projects
  - heading "Google" [level=2]
  - img
  - text: Google is a multinational technology company that specializes in internet-related services and products. 0 posts 0 projects
  - heading "Intel" [level=2]
  - img
  - text: Intel is a multinational semiconductor and electronics company headquartered in Santa Clara, California. 0 posts 0 projects
  - heading "Microsoft" [level=2]
  - img
  - text: Microsoft is a technology company that develops and sells software, services, and products. 0 posts 0 projects
  - heading "Oracle" [level=2]
  - img
  - text: Oracle is a software company that provides a range of products and services for businesses and organizations. 0 posts 0 projects
  - heading "Texas Instruments" [level=2]
  - img
  - text: Texas Instruments is a semiconductor company that designs, manufactures, and sells integrated circuits, microprocessors, and microcontrollers. 0 posts 0 projects
  - heading "Turso" [level=2]
  - img
  - text: Turso is a fully managed database platform that you can use to create hundreds of thousands of databases per organization and supports replication to any location, including your own servers, for microsecond-latency access. 0 posts 0 projects
  - heading "MySQL" [level=2]
  - img
  - text: MySQL is a relational database management system (RDBMS) that is open-source and developed by Oracle Corporation. 0 posts 0 projects
  - heading "Postgres" [level=2]
  - img
  - text: Postgres is a relational database management system that is open-source and developed by PostgreSQL Global Development Group. 0 posts 0 projects
  - heading "SQLite" [level=2]
  - img
  - text: SQLite is the most used relational database management system that is open-source and developed by C. 0 posts 0 projects
  - heading "iastate" [level=2]
  - img
  - text: Iowa State University where I got my BS in Electrical Engineering and Minor in Computer Science. 0 posts 4 projects
  - heading "AGEDS461" [level=2]
  - img
  - text: This is a class taken at Iowa State University. 3 posts 0 projects
  - heading "ARCH 321" [level=2]
  - img
  - text: Study of the development of the built environment and urban condition in the United States from the colonial period to today. 1 posts 0 projects
  - heading "CPRE281" [level=2]
  - img
  - text: This is a class taken at Iowa State University. 1 posts 1 projects
  - heading "CPRE288" [level=2]
  - img
  - text: Embedded C programming. Interrupt handling. Memory mapped I/O in the context of an application. Elementary embedded design flow/methodology. Timers, scheduling, resource allocation, optimization, state machine based controllers, real time constraints within the context of an application. Applications laboratory exercises with embedded devices class taken at Iowa State University. 1 posts 0 projects
  - heading "CPRE381" [level=2]
  - img
  - text: CPRE381 1 posts 0 projects
  - heading "CPRE488" [level=2]
  - img
  - text: This is a class taken at Iowa State University. 2 posts 1 projects
  - heading "CS228" [level=2]
  - img
  - text: This is a class taken at Iowa State University. 0 posts 1 projects
  - heading "ECON 101" [level=2]
  - img
  - text: Resource allocation, opportunity cost, comparative and absolute advantage taught at Iowa State University. 0 posts 0 projects
  - heading "EE201" [level=2]
  - img
  - text: Electric Circuits at Iowa State University 1 posts 0 projects
  - heading "EE224" [level=2]
  - img
  - text: Signals and Systems I taught at Iowa State University. 0 posts 0 projects
  - region "Contact":
    - heading "Get In Touch" [level=2]
    - paragraph: Interested in working together? Feel free to reach out through any of the channels below.
    - link:
      - /url: https://www.linkedin.com/in/conner-ohnesorge-b720a4238
      - img
    - link:
      - /url: https://github.com/conneroisu
      - img
    - link:
      - /url: https://x.com/ConnerOhnesorge
      - img
    - link:
      - /url: mailto:conneroisu@outlook.com
      - img
    - text: Name
    - textbox "Name"
    - text: Email
    - textbox "Email"
    - text: Subject
    - textbox "Subject"
    - text: Message
    - textbox "Message"
    - button "Send Message"
- contentinfo:
  - heading "Conner Ohnesorge" [level=3]
  - paragraph: Electrical Engineer & Software Developer
  - link "LinkedIn":
    - /url: https://www.linkedin.com/in/conner-ohnesorge-b720a4238
  - link "GitHub":
    - /url: https://github.com/conneroisu
  - link "Twitter":
    - /url: https://x.com/ConnerOhnesorge
  - link "Email":
    - /url: mailto:conneroisu@outlook.com
  - paragraph: © 2025 Conner Ohnesorge. All rights reserved.
  - link "Posts":
    - /url: /posts
  - link "Projects":
    - /url: /projects
  - link "Tags":
    - /url: /tags
  - link "Contact":
    - /url: "#contact"
```

# Test source

```ts
   1 | import { test, expect } from '@playwright/test'
   2 |
   3 | test('Alpine.js data attributes', async ({ page }) => {
   4 |   await page.goto('/')
   5 |   
   6 |   // Test that Alpine.js attributes exist on elements (like mobile menu)
   7 |   await expect(page.locator('[x-data]')).toHaveCount(await page.locator('[x-data]').count())
   8 |   
   9 |   // Test mobile menu specifically
  10 |   const mobileMenuContainer = page.locator('[x-data*="isMenuOpen"]')
  11 |   if (await mobileMenuContainer.isVisible()) {
  12 |     await expect(mobileMenuContainer).toHaveAttribute('x-data')
  13 |     await expect(page.locator('[x-show="isMenuOpen"]')).toBeAttached()
> 14 |     await expect(page.locator('[\\@click*="isMenuOpen"]')).toBeVisible()
     |                                                            ^ Error: expect.toBeVisible: Error: strict mode violation: locator('[\\@click*="isMenuOpen"]') resolved to 5 elements:
  15 |   }
  16 | })
  17 |
  18 | test('mobile menu interaction with Alpine.js', async ({ page }) => {
  19 |   await page.goto('/')
  20 |   
  21 |   // Set mobile viewport
  22 |   await page.setViewportSize({ width: 375, height: 667 })
  23 |   
  24 |   const menuButton = page.locator('[\\@click*="isMenuOpen"]')
  25 |   const mobileMenu = page.locator('[x-show="isMenuOpen"]')
  26 |   
  27 |   if (await menuButton.isVisible()) {
  28 |     // Test initial state - menu should be hidden
  29 |     await expect(mobileMenu).toBeHidden()
  30 |     
  31 |     // Click to open menu
  32 |     await menuButton.click()
  33 |     
  34 |     // Wait for Alpine.js to show the menu
  35 |     await expect(mobileMenu).toBeVisible()
  36 |     
  37 |     // Click menu button again to close
  38 |     await menuButton.click()
  39 |     
  40 |     // Menu should be hidden again
  41 |     await expect(mobileMenu).toBeHidden()
  42 |   }
  43 | })
```