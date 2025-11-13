import { Link } from "@tanstack/solid-router";

export default function Footer() {
	const currentYear = new Date().getFullYear();

	return (
		<footer class="bg-gray-800 border-t border-gray-700 py-12">
			<div class="container mx-auto px-4">
				{/* Top Section - Name and Social Links */}
				<div class="flex flex-col md:flex-row justify-between items-center">
					{/* Left Side - Name and Title */}
					<div class="mb-6 md:mb-0">
						<h3 class="text-white text-xl font-bold mb-2">Conner Ohnesorge</h3>
						<p class="text-gray-400">
							Electrical Engineer & Software Developer
						</p>
					</div>

					{/* Right Side - Social Links */}
					<div class="flex flex-wrap gap-4 justify-center">
						<a
							href="https://www.linkedin.com/in/conner-ohnesorge-b720a4238"
							target="_blank"
							rel="noopener noreferrer"
							class="text-gray-400 hover:text-green-400 transition-colors"
						>
							LinkedIn
						</a>
						<a
							href="https://github.com/conneroisu"
							target="_blank"
							rel="noopener noreferrer"
							class="text-gray-400 hover:text-green-400 transition-colors"
						>
							GitHub
						</a>
						<a
							href="https://x.com/ConnerOhnesorge"
							target="_blank"
							rel="noopener noreferrer"
							class="text-gray-400 hover:text-green-400 transition-colors"
						>
							Twitter
						</a>
						<a
							href="mailto:conneroisu@outlook.com"
							class="text-gray-400 hover:text-green-400 transition-colors"
						>
							Email
						</a>
					</div>
				</div>

				{/* Bottom Section - Copyright and Footer Navigation */}
				<div class="mt-8 pt-8 border-t border-gray-700 flex flex-col md:flex-row justify-between items-center">
					{/* Copyright */}
					<p class="text-gray-500 text-sm">
						&copy; {currentYear} Conner Ohnesorge. All rights reserved.
					</p>

					{/* Footer Navigation Links */}
					<div class="mt-4 md:mt-0">
						<Link
							to="/posts"
							class="text-gray-500 hover:text-gray-300 transition-colors text-sm mx-2"
						>
							Posts
						</Link>
						<Link
							to="/projects"
							class="text-gray-500 hover:text-gray-300 transition-colors text-sm mx-2"
						>
							Projects
						</Link>
						<Link
							to="/tags"
							class="text-gray-500 hover:text-gray-300 transition-colors text-sm mx-2"
						>
							Tags
						</Link>
						<a
							href="#contact"
							class="text-gray-500 hover:text-gray-300 transition-colors text-sm mx-2"
						>
							Contact
						</a>
					</div>
				</div>
			</div>
		</footer>
	);
}
