INSERT
    OR IGNORE INTO posts (
        title,
        description,
        slug,
        content,
        banner_url,
        created_at,
        updated_at
    )
VALUES
    (
        'Getting Started with Go',
        'Learn how to set up your Go development environment and create your first application.',
        'getting-started-with-go',
        '# Getting Started with Go Congratulations! You have written and executed your first Go program.',
        '/assets/img/posts/golang.jpg',
        '1577836800',
        '1577836800'
    ),
    (
        'Introduction to SQLite',
        'Learn how to use SQLite for simple database applications.',
        'introduction-to-sqlite',
        '# Introduction to SQLite Go has excellent support for SQLite through the database/sql package and drivers like modernc.org/sqlite.',
        '/assets/img/posts/sqlite.jpg',
        '1583020800',
        '1583020800'
    ),
    (
        'Building Web Applications with Templ',
        'Learn how to use the Templ library to create type-safe HTML templates in Go.',
        'building-with-templ',
        '# Building Web Applications with Templ Templ makes it simple to create reusable UI components that seamlessly integrate with your Go code.',
        '/assets/img/posts/templates.jpg',
        '1588291200',
        '1588291200'
    ),
    (
        'Modern Frontend with HTMX and Alpine.js',
        'Simplify your frontend architecture by combining HTMX for AJAX and Alpine.js for interactivity.',
        'modern-frontend-htmx-alpine',
        '# Modern Frontend with HTMX and Alpine.js HTMX and Alpine.js work wonderfully together. HTMX handles server communication and DOM updates, while Alpine handles client-side interactivity and state.',
        '/assets/img/posts/frontend.jpg',
        '1593561600',
        '1593561600'
    ),
    (
        'Modern Frontend with Tailwind CSS',
        'Tailwind CSS is a utility-first CSS framework that allows you to build highly customizable and responsive designs.',
        'modern-frontend-tailwind',
        '# Modern Frontend with Tailwind CSS Tailwind CSS is a utility-first CSS framework that allows you to build highly customizable and responsive designs. It provides a set of predefined classes that you can use to quickly style your HTML elements.',
        '/assets/img/posts/frontend.jpg',
        '1593561600',
        '1593561600'
    ),
    (
        'NixOS Development Environment',
        'Learn how to set up a NixOS development environment and build your first NixOS application.',
        'nixos-development-environment',
        '# NixOS Development Environment NixOS is a powerful and flexible operating system that allows you to build your own custom Linux distribution.',
        '/assets/img/posts/nixos.jpg',
        '1593561600',
        '1593561600'
    );

-- Let's also create some associations between posts and tags/projects
INSERT INTO
    post_tags (post_id, tag_id)
VALUES
    (1, 1),
    -- "Getting Started with Go" with "Tag 1"
    (1, 2),
    -- "Getting Started with Go" with "Tag 2"
    (2, 2),
    -- "Introduction to SQLite" with "Tag 2"
    (3, 3),
    -- "Building Web Applications with Templ" with "Tag 3"
    (4, 1),
    -- "Modern Frontend with HTMX and Alpine.js" with "Tag 1"
    (5, 3);

-- "NixOS Development Environment" with "Tag 3"
INSERT INTO
    post_projects (post_id, project_id)
VALUES
    (1, 1),
    -- "Getting Started with Go" with "Project 1"
    (2, 1),
    -- "Introduction to SQLite" with "Project 1"
    (3, 2),
    -- "Building Web Applications with Templ" with "Project 2"
    (4, 3),
    -- "Modern Frontend with HTMX and Alpine.js" with "Project 3"
    (5, 4);

-- "NixOS Development Environment" with "Project 4"
INSERT INTO
    project_tags (project_id, tag_id)
VALUES
    (4, 4);

-- "NixOS Development Environment" with "Tag 4"
