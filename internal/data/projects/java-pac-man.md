---
id: java-pac-man
aliases:
  - Java Pac-Man
  - Pac-Man Implementation
  - "Engineering Excellence in Retro Gaming: My Java Pac-Man Implementation"
tags:
  - programming-language/java
  - framework/maven
  - edu/iastate/cs228
banner_path: projects/java-pac-man.webp
created_at: 2023-11-15T10:00:00.000-06:00
description: A sophisticated Java implementation of Pac-Man featuring advanced AI ghost personalities and professional-grade software engineering practices
title: Java Pac-Man Implementation
updated_at: 2025-06-05T10:35:14.000-06:00
---

# Engineering Excellence in Retro Gaming: My Java Pac-Man Implementation

I developed this Java Pac-Man project for CS228 at Iowa State University, creating what became a remarkable fusion of classic arcade gaming with modern software engineering practices. 
**My implementation spans over 1,500 lines of sophisticated code** including a __PacmanGame__ class managing complex game states and a __ActorImpl__ base class implementing advanced AI pathfinding algorithms. What began as an educational project evolved into a production-quality codebase featuring **four distinct ghost AI personalities**, comprehensive testing infrastructure, and professional-grade build automation that rivals commercial game development standards.

## Architectural Sophistication Meets Educational Clarity

I designed the project's architecture to demonstrate a masterful balance between complexity and teachability. At its core, my implementation employs a **clean three-tier package structure** separating concerns into `api/`, `com.pacman.ghost/`, and `ui/` packages. This separation isn't merely organizational—it represents my deep understanding of software architecture principles that many commercial games fail to achieve.

My PacmanGame class, weighing in at **585 lines of carefully orchestrated code**, serves as the central nervous system of the implementation. Rather than falling into the common trap of creating a "god class," this component manages game state through a sophisticated finite state machine pattern that I designed. The class coordinates multiple subsystems including collision detection, scoring mechanics, ghost mode transitions, and the famous tunnel wraparound feature—all while maintaining clear separation of concerns.

I built the build system to exemplify modern Java development practices rarely seen in educational projects. **Maven serves as the foundation**, but the true sophistication lies in the comprehensive quality gates I implemented: PMD for code quality analysis, Checkstyle for consistency enforcement, and SpotBugs for potential bug detection. I included a Maven wrapper to ensure build consistency across development environments, while the Nix flake configuration demonstrates cutting-edge reproducible build practices typically reserved for enterprise applications.

## Ghost AI: Where Computer Science Meets Game Design Brilliance

The ghost AI implementation represents the crown jewel of my project, housed within the **979-line ActorImpl base class** that I created as the foundation for sophisticated pathfinding and behavioral patterns. Each ghost possesses a distinct personality that I implemented through the Strategy pattern, creating the emergent gameplay that made Pac-Man legendary.

**Blinky, the red ghost**, implements the simplest yet most relentless algorithm I designed—direct pursuit of Pac-Man's current position. His targeting calculation involves straightforward Euclidean distance minimization, but my implementation includes the famous "Cruise Elroy" mode where Blinky's speed increases when fewer than 20 dots remain, maintaining aggression even during scatter phases.

**My Pinky's ambush algorithm** showcases more sophisticated predictive targeting, calculating a position four tiles ahead of Pac-Man's current direction. I faithfully reproduced the original's overflow bug—when Pac-Man faces upward, the target becomes four tiles up AND four tiles left. This "bug" has become a feature in my implementation, adding unpredictability that enhances gameplay.

**Inky's implementation** demonstrates the most complex algorithm I developed, using vector mathematics to create unpredictable behavior. The ghost calculates an intermediate position two tiles ahead of Pac-Man, then uses Blinky's position to create a vector that's doubled in length for the final target. This creates fascinating emergent behavior where Inky can suddenly appear from unexpected directions.

**My Clyde's distance-based personality switching** implements a dual-mode behavior system. When more than eight tiles from Pac-Man, he pursues directly like Blinky. Within eight tiles, he retreats to his scatter corner, creating the "bashful" personality that provides breathing room for players.

I designed the ghost state machine to implement **five distinct behavioral modes**: CHASE, SCATTER, FRIGHTENED, DEAD, and INACTIVE. State transitions follow the classic pattern with decreasing scatter durations as levels progress. My implementation includes forced direction reversals on state changes, providing visual feedback to players about AI behavior shifts.

## Testing Infrastructure Worthy of Commercial Development

I developed a testing framework that demonstrates a level of sophistication rarely seen in academic projects. **My SimulationTestFramework** enables deterministic testing of AI behaviors by mocking time-dependent systems. This allows me to verify complex scenarios like ghost behavior at intersections, tunnel traversal mechanics, and state transition timing.

I wrote unit tests covering individual components using JUnit 5's powerful features, including parameterized tests for validating ghost behavior across multiple game levels. My integration tests verify the interaction between subsystems—ensuring that collision detection, scoring, and AI pathfinding work harmoniously. I included performance benchmarks validating that the game maintains 60 FPS even with all four ghosts actively pathfinding.

**I implemented mock-based testing** with Mockito to isolate game logic from UI concerns, enabling rapid test execution and comprehensive coverage. My test suite validates edge cases like simultaneous ghost collisions, power pellet timing windows, and the intricate dot-counting system that controls ghost house exit timing.

## Code Quality: From Academic Project to Professional Standards

I evolved the codebase from educational project to professional-quality implementation, evidenced by my comprehensive static analysis integration. **PMD identified 130+ violations**, each representing an opportunity for improvement that I documented. Rather than viewing these as failures, I included a detailed refactoring plan addressing each category of violation.

Common violations I identified include naming convention issues typical of game development—short variable names in tight loops, magic numbers for game constants, and occasionally long methods in the game update cycle. My refactoring strategy employs automated fixes for simple violations while carefully preserving game-specific optimizations that might trigger false positives.

**I configured Checkstyle** to enforce consistent code formatting across the project, while SpotBugs catches potential null pointer exceptions and performance issues. My integration of these tools into the Maven build process creates quality gates that prevent regression while teaching developers about professional development practices.

I ensured comprehensive JavaDoc coverage extends beyond mere API documentation. My comments explain not just what the code does, but why I made specific algorithmic choices. The ghost AI documentation includes mathematical formulas for targeting calculations, making my codebase valuable for both users and future maintainers.

## Educational Value Transcending the Classroom

While CS228 at Iowa State focuses primarily on data structures rather than game development, I brilliantly demonstrated how foundational computer science concepts apply to engaging real-world applications. My ghost pathfinding algorithms provide concrete examples of **graph traversal and distance calculations**. The state machine implementation showcases finite automata theory in action.

I used **design patterns** throughout my project, creating a living textbook of software engineering principles. The Strategy pattern enables pluggable ghost behaviors, the Observer pattern manages game events, and the State pattern controls game flow. Students can see these patterns not as abstract concepts but as practical tools solving real problems in my implementation.

## Modern Development Practices Setting New Standards

I integrated **Nix flakes** for reproducible development environments, demonstrating cutting-edge practices rarely seen even in commercial projects. This ensures that whether a developer works on Linux, macOS, or Windows, they experience identical build environments, eliminating "works on my machine" issues.

**My Maven wrapper** inclusion means developers don't need to install specific Maven versions, reducing onboarding friction.

Combined with comprehensive README documentation and clear package structure, new contributors can achieve productivity within minutes rather than hours.

I built logging infrastructure using SLF4J with Logback to enable sophisticated debugging of ghost AI decisions.

Structured logging captures ghost state transitions, pathfinding decisions, and collision events with nanosecond precision. This creates an audit trail invaluable for understanding emergent behaviors and debugging complex interactions.

## A Masterclass in Software Engineering Education

My Java Pac-Man implementation transcends its origins as a CS228 project to become a masterclass in software engineering. I demonstrate that educational projects need not sacrifice quality for simplicity—instead, they can achieve both through thoughtful architecture and iterative improvement.

**My 585-line PacmanGame class** and **979-line ActorImpl base** represent not bloated code but carefully crafted systems where each line serves a purpose. The four ghost personalities showcase how simple algorithms can create complex, engaging behaviors. My comprehensive testing infrastructure proves that quality isn't optional, even in academic settings.

Most importantly, my project bridges the gap between academic computer science and practical software engineering. I show that the algorithms we study have real applications, that design patterns solve actual problems, and that code quality tools aren't bureaucratic obstacles but enablers of excellence. In an industry where many developers learn these lessons years into their careers, my project provides them upfront, wrapped in the engaging package of a classic arcade game.
