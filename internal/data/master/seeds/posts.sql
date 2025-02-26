INSERT INTO
    posts (
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
    );

INSERT INTO
    posts (
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
        'Introduction to SQLite Databases',
        'A comprehensive guide to getting started with SQLite for simple database applications.',
        'introduction-to-sqlite',
        '# Introduction to SQLite Go has excellent support for SQLite through the database/sql package and drivers like modernc.org/sqlite.',
        '/assets/img/posts/sqlite.jpg',
        '1583020800',
        '1583020800'
    );

INSERT INTO
    posts (
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
        'Building Web Applications with Templ',
        'Learn how to use the Templ library to create type-safe HTML templates in Go.',
        'building-with-templ',
        '# Building Web Applications with Templ Templ makes it simple to create reusable UI components that seamlessly integrate with your Go code.',
        '/assets/img/posts/templates.jpg',
        '1588291200',
        '1588291200'
    );

INSERT INTO
    posts (
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
        'Modern Frontend with HTMX and Alpine.js',
        'Simplify your frontend architecture by combining HTMX for AJAX and Alpine.js for interactivity.',
        'modern-frontend-htmx-alpine',
        '# Modern Frontend with HTMX and Alpine.js HTMX and Alpine.js work wonderfully together. HTMX handles server communication and DOM updates, while Alpine handles client-side interactivity and state.',
        '/assets/img/posts/frontend.jpg',
        '1593561600',
        '1593561600'
    );

INSERT INTO
    posts (
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
        'NixOS Development Environment',
        'How to set up a reproducible development environment using Nix and flakes.',
        'nixos-development-environment',
        '# NixOS Development Environment This command will drop you into a shell with all the specified dependencies available.',
        '/assets/img/posts/nixos.jpg',
        '1598832000',
        '1598832000'
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
