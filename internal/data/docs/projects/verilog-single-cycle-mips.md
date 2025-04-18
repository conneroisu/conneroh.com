---
id: verilog-single-cycle-mips
aliases:
  - singe-cycle-mips-verilog
tags:
  - hdl/verilog
  - edu/iastate/cpre281
banner_path: projects/single-cycle-mips-verilog.webp
created_at: 2025-04-12T19:35:20.000-06:00
description: A single cycle mips processor implementation in verilog.
title: Verilog Single Cycle MIPS
updated_at: 2025-04-15T05:32:18.000-06:00
---

Github: [https://github.com/conneroisu/single-cycle-mips-verilog](https://github.com/conneroisu/single-cycle-mips-verilog)

# singe-cycle-mips-verilog

A single cycle mips processor implementation in verilog with reflections from writing a similar single-cycle
processor project in vhdl that can execute a subset of the MIPS
instruction set and displays the current instruction on the seven segment displays present
on the FPGA board.
The processor is implemented in Verilog and tested using a
test-bench to verify the functionality of the processor. Additionally, the processor
is tested on an FPGA board to verify the functionality of the processor.

<body>
</div>
</header>
<nav id="TOC" role="doc-toc">
<h2 id="toc-title">Table of Contents</h2>
<ul>
<li><a href="#proposal" id="toc-proposal">Proposal</a></li>
<li><a href="#introduction" id="toc-introduction">Introduction</a></li>
<li><a href="#state-machines" id="toc-state-machines">State Machines</a>
<ul>
<li><a href="#single-cycle-mips-processor-staging"
id="toc-single-cycle-mips-processor-staging">Single Cycle MIPS Processor
Staging</a></li>
<li><a href="#execution-diagram-for-each-instruction"
id="toc-execution-diagram-for-each-instruction">Execution diagram for
each instruction:</a></li>
<li><a href="#add-r-type-instruction"
id="toc-add-r-type-instruction">ADD (R-type instruction)</a></li>
<li><a href="#addi-i-type-instruction"
id="toc-addi-i-type-instruction">ADDI (I-type instruction)</a></li>
<li><a href="#lw-load-word" id="toc-lw-load-word">LW (Load
Word)</a></li>
<li><a href="#beq-branch-if-equal" id="toc-beq-branch-if-equal">BEQ
(Branch if Equal)</a></li>
<li><a href="#j-jump" id="toc-j-jump">J (Jump)</a></li>
<li><a href="#addu-add-unsigned" id="toc-addu-add-unsigned">ADDU (Add
Unsigned)</a></li>
<li><a href="#sub-subtract" id="toc-sub-subtract">SUB
(Subtract)</a></li>
<li><a href="#subu-subtract-unsigned"
id="toc-subu-subtract-unsigned">SUBU (Subtract Unsigned)</a></li>
<li><a href="#and-bitwise-and" id="toc-and-bitwise-and">AND (Bitwise
AND)</a></li>
<li><a href="#andi-and-immediate" id="toc-andi-and-immediate">ANDI (AND
Immediate)</a></li>
<li><a href="#or" id="toc-or">OR</a></li>
<li><a href="#ori-or-immediate" id="toc-ori-or-immediate">ORI (OR
Immediate)</a></li>
<li><a href="#nor" id="toc-nor">NOR</a></li>
</ul></li>
<li><a href="#comparing-verilog-vs-vhdl"
id="toc-comparing-verilog-vs-vhdl">Comparing Verilog vs VHDL</a>
<ul>
<li><a href="#interesting-notes-about-verilog"
id="toc-interesting-notes-about-verilog">Interesting notes about
verilog</a></li>
<li><a href="#interesting-notes-about-vhdl"
id="toc-interesting-notes-about-vhdl">Interesting notes about
VHDL</a></li>
<li><a href="#conclusion" id="toc-conclusion">Conclusion</a></li>
</ul></li>
<li><a href="#breaking-down-decoding-a-signal-to-7-segment-displays"
id="toc-breaking-down-decoding-a-signal-to-7-segment-displays">Breaking
down decoding a signal to 7-segment displays</a>
<ul>
<li><a href="#schematics" id="toc-schematics">Schematics</a>
<ul>
<li><a href="#control-unit-schematic"
id="toc-control-unit-schematic">Control Unit Schematic</a></li>
<li><a href="#register-file" id="toc-register-file">Register
File</a></li>
<li><a href="#data-memory" id="toc-data-memory">Data Memory</a></li>
<li><a href="#alu-control" id="toc-alu-control">ALU Control</a></li>
<li><a href="#program-counter-control"
id="toc-program-counter-control">Program Counter Control</a></li>
<li><a href="#alu" id="toc-alu">ALU</a></li>
<li><a href="#instruction-memory"
id="toc-instruction-memory">Instruction Memory</a></li>
<li><a href="#program-counter" id="toc-program-counter">Program
Counter</a></li>
<li><a href="#waveform" id="toc-waveform">Waveform</a></li>
</ul></li>
</ul></li>
<li><a href="#tooling" id="toc-tooling">Tooling</a>
<ul>
<li><a href="#components-and-explanations"
id="toc-components-and-explanations">Components and
Explanations</a></li>
<li><a href="#verbose-components-code"
id="toc-verbose-components-code">Verbose Components Code</a>
<ul>
<li><a href="#data-memory-2" id="toc-data-memory-2">Data Memory</a></li>
<li><a href="#instruction-memory-2"
id="toc-instruction-memory-2">Instruction Memory</a></li>
<li><a href="#program-counter-1" id="toc-program-counter-1">Program
Counter</a></li>
<li><a href="#alu-2" id="toc-alu-2">ALU</a></li>
<li><a href="#control-unit-1" id="toc-control-unit-1">Control
Unit</a></li>
<li><a href="#testbench" id="toc-testbench">Testbench</a></li>
</ul></li>
</ul></li>
<li><a href="#conclusion-1" id="toc-conclusion-1">Conclusion</a></li>
</ul>
</nav>
<p>The processor is able to execute the
following instructions: <strong>LW</strong> <strong>SW</strong>
<strong>J</strong> <strong>ADD</strong> <strong>ADDI</strong>
<strong>BEQ</strong> <strong>ADDU</strong> <strong>SUBU</strong>
<strong>AND</strong> <strong>ANDI</strong> <strong>OR</strong>
<strong>ORI</strong> <strong>SUB</strong> <strong>NOR</strong>
<strong>BNE</strong> <strong>SLT</strong></p>
<p>
The current instruction being executed will be displayed
on the seven segment displays present on the FPGA board. The processor
will be implemented in Verilog and tested using a test-bench.</p>
<p>After the processor has been verified using the test-bench, the
processor will be tested on an FPGA board to verify the functionality of
the processor.</p>
<p>The comparison and contrast of the experience writing the same
processor in both Verilog and VHDL will be also included in the final
report.</p>
<h2 id="introduction">Introduction</h2>
<p>The project is a single-cycle MIPS processor that can at a variable
speed execute a subset of the MIPS instruction set displaying the
current instruction on the seven segment displays present on
<strong>EP4CE115F29C7</strong> FPGA board.</p>
<div
style="max-width: fit-content; margin-left: auto; margin-right: auto;">
<figure>
<img src="assets/Circuit%20Diagram.png" alt="Circuit Diagram.png" />
<figcaption aria-hidden="true">Circuit Diagram.png</figcaption>
</figure>
</div>
<p>As a result of this desired function (and it’s actualation into
reality), the processor can technically be used to do
<strong>all</strong> the other projects from other students in the
class. As I took this class whilst also taking CPRE381, I additionally
decided to compare and contrast the experience writing the same
processor in both Verilog and VHDL.</p>
<p>While the overall state machine will be broken down below, the main
processor state machine has <strong>five</strong> states:</p>
<p><strong>Fetch</strong>: In this state, the processor fetches the next
instruction from memory. <strong>Decode</strong>: In this state, the
processor decodes the instruction to determine what operation to
perform. <strong>Execute</strong>: In this state, the processor executes
the instruction. <strong>Memory</strong>: In this state, the processor
accesses memory to read or write data. <strong>Write-back</strong>: In
this state, the processor writes the results of the instruction to a
register.</p>
<p>Supported Instructions: LW SW J ADD ADDI BEQ ADDU SUBU AND ANDI OR
ORI SUB NOR BNE</p>
<h2 id="state-machines">State Machines</h2>

![Assembly-Instructions-Architecture.png](assets/Assembly-Instructions-Architecture.png)

<p>Each instruction type (R-type, I-type, J-type) generally follows a
similar flow with variations primarily in the Execute and Memory Access
stages depending on whether the instruction involves arithmetic, memory
access, or control flow.</p>
<p>This model provides a consistent framework for understanding how


different instructions are processed in the single-cycle MIPS
architecture.</p>


![Single-Cycle-Staging.png](assets/Single-Cycle-Staging.png)

<h3 id="add-r-type-instruction">ADD (R-type instruction)</h3>
<p>Name (format, op, function): <code>add (R,0,32)</code></p>
<p>Syntax: <code>add rd,rs,rt</code></p>
<p>Operation: <code>reg(rd) := reg(rs) + reg(rt);</code></p>
<h4 id="instruction-overview">Instruction Overview:</h4>
<p>The following is an overview of the operation of the <code>add</code>
instruction in a MIPS processor, focusing on the key stages of the
processor pipeline.</p>
<ul>
<li><strong>IF:</strong> The instruction is fetched from memory using
the program counter (PC).</li>
<li><strong>ID:</strong> The instruction bits are decoded to determine
it is an ADD operation. Registers specified by the source register
fields (<code>rs</code> and <code>rt</code>) are read.</li>
<li><strong>EX:</strong> The ALU performs the addition of the two
register values.</li>
<li><strong>MEM:</strong> No action (not used by ADD).</li>
<li><strong>WB:</strong> The result from the ALU is written back to the
destination register (<code>rd</code>).</li>
</ul>
<h4 id="operation-breakdown">Operation Breakdown:</h4>
<p>The following provides a detailed breakdown of the operation of the
<code>add</code> instruction in a MIPS processor, focusing on the key
stages of the processor pipeline.</p>
<h5 id="stages-of-the-add-instruction">Stages of the <code>add</code>
Instruction:</h5>
<ol type="1">
<li><strong>Instruction Fetch (IF):</strong>
<ul>
<li><p>The instruction is fetched from memory using the Program Counter
(PC).</p></li>
<li><p>This stage corresponds to reading the instruction code from the
instruction memory. The address comes from the PC which points to the
location of the next instruction to execute.</p>
<div class="sourceCode" id="cb1"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb1-1"><a href="#cb1-1" aria-hidden="true" tabindex="-1"></a>i_Instruction <span class="op">=</span> Imem<span class="op">[</span>i_Addr<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span></code></pre></div></li>
</ul></li>
<li><strong>Instruction Decode (ID):</strong>
<ul>
<li><p>The fetched instruction is decoded to determine it is a
<code>add</code> operation.</p></li>
<li><p>The opcode part of the instruction (which is <code>000000</code>
for R-type instructions) is identified, and the source register
identifiers (<code>rs</code> and <code>rt</code>) are used to read the
respective registers.</p>
<div class="sourceCode" id="cb2"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb2-1"><a href="#cb2-1" aria-hidden="true" tabindex="-1"></a>rs <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb2-2"><a href="#cb2-2" aria-hidden="true" tabindex="-1"></a>rt <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span></code></pre></div></li>
</ul></li>
<li><strong>Execute (EX):</strong>
<ul>
<li><p>The ALU (Arithmetic Logic Unit) performs the addition of the
values in the source registers (<code>rs</code> and
<code>rt</code>).</p></li>
<li><p>The values from these registers are fed into the ALU where the
addition is performed based on the control signal (<code>ALUOp</code>)
from the control unit.</p>
<div class="sourceCode" id="cb3"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb3-1"><a href="#cb3-1" aria-hidden="true" tabindex="-1"></a>o_ALUresult <span class="op">=</span> i_data1 <span class="op">+</span> i_data2<span class="op">;</span></span></code></pre></div></li>
</ul></li>
<li><strong>Memory Access (MEM):</strong>
<ul>
<li><p>For the <code>add</code> instruction, this stage is not utilized
as no memory access is required (i.e., no data is read from or written
to the memory).</p>
<div class="sourceCode" id="cb4"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb4-1"><a href="#cb4-1" aria-hidden="true" tabindex="-1"></a><span class="co">// No memory operation required for &#39;add&#39;</span></span></code></pre></div></li>
</ul></li>
<li><strong>Write Back (WB):</strong>
<ul>
<li><p>The result from the ALU is written back into the destination
register (<code>rd</code>).</p></li>
<li><p>This is where the output of the ALU operation is stored back into
the register file, specifically into the register indicated by the
<code>rd</code> field of the instruction.</p>
<div class="sourceCode" id="cb5"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb5-1"><a href="#cb5-1" aria-hidden="true" tabindex="-1"></a>Reg<span class="op">[</span>rd<span class="op">]</span> <span class="op">=</span> o_ALUresult<span class="op">;</span></span></code></pre></div></li>
</ul></li>
</ol>
<h5 id="example-code-snippet">Example Code Snippet:</h5>
<p>Here is a simplified Verilog snippet that captures the essence of the
<code>add</code> instruction’s operation in a MIPS processor, focusing
on the key stages of the mips architecture (IF, ID, EX, MEM, WB).</p>
<div class="sourceCode" id="cb6"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb6-1"><a href="#cb6-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> MIPS_Processor<span class="op">(</span><span class="dt">input</span> clk<span class="op">,</span> <span class="dt">input</span> reset<span class="op">,</span> ...<span class="op">);</span></span>
<span id="cb6-2"><a href="#cb6-2" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Registers and other declarations here</span></span>
<span id="cb6-3"><a href="#cb6-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> PC<span class="op">,</span> ALUResult<span class="op">,</span> Reg<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb6-4"><a href="#cb6-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> InstructionRegister<span class="op">,</span> ReadData1<span class="op">,</span> ReadData2<span class="op">;</span></span>
<span id="cb6-5"><a href="#cb6-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">integer</span> rd<span class="op">,</span> rs<span class="op">,</span> rt<span class="op">;</span></span>
<span id="cb6-6"><a href="#cb6-6" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb6-7"><a href="#cb6-7" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>reset<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb6-8"><a href="#cb6-8" aria-hidden="true" tabindex="-1"></a>            PC <span class="op">&lt;=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// Reset PC</span></span>
<span id="cb6-9"><a href="#cb6-9" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb6-10"><a href="#cb6-10" aria-hidden="true" tabindex="-1"></a>            <span class="co">// Fetch Instruction</span></span>
<span id="cb6-11"><a href="#cb6-11" aria-hidden="true" tabindex="-1"></a>            InstructionRegister <span class="op">&lt;=</span> Imem<span class="op">[</span>PC<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span>
<span id="cb6-12"><a href="#cb6-12" aria-hidden="true" tabindex="-1"></a>            PC <span class="op">&lt;=</span> PC <span class="op">+</span> <span class="dv">4</span><span class="op">;</span></span>
<span id="cb6-13"><a href="#cb6-13" aria-hidden="true" tabindex="-1"></a>            <span class="co">// Decode Instruction</span></span>
<span id="cb6-14"><a href="#cb6-14" aria-hidden="true" tabindex="-1"></a>            rs <span class="op">=</span> InstructionRegister<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span> <span class="co">// Extract source register indices</span></span>
<span id="cb6-15"><a href="#cb6-15" aria-hidden="true" tabindex="-1"></a>            rt <span class="op">=</span> InstructionRegister<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span> <span class="co">// Extract target register indices</span></span>
<span id="cb6-16"><a href="#cb6-16" aria-hidden="true" tabindex="-1"></a>            rd <span class="op">=</span> InstructionRegister<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">11</span><span class="op">];</span> <span class="co">// Extract destination register index</span></span>
<span id="cb6-17"><a href="#cb6-17" aria-hidden="true" tabindex="-1"></a>            ReadData1 <span class="op">&lt;=</span> Reg<span class="op">[</span>rs<span class="op">];</span></span>
<span id="cb6-18"><a href="#cb6-18" aria-hidden="true" tabindex="-1"></a>            ReadData2 <span class="op">&lt;=</span> Reg<span class="op">[</span>rt<span class="op">];</span></span>
<span id="cb6-19"><a href="#cb6-19" aria-hidden="true" tabindex="-1"></a>            <span class="co">// Execute</span></span>
<span id="cb6-20"><a href="#cb6-20" aria-hidden="true" tabindex="-1"></a>            ALUResult <span class="op">&lt;=</span> ReadData1 <span class="op">+</span> ReadData2<span class="op">;</span></span>
<span id="cb6-21"><a href="#cb6-21" aria-hidden="true" tabindex="-1"></a>            <span class="co">// Write Back</span></span>
<span id="cb6-22"><a href="#cb6-22" aria-hidden="true" tabindex="-1"></a>            Reg<span class="op">[</span>rd<span class="op">]</span> <span class="op">&lt;=</span> ALUResult<span class="op">;</span></span>
<span id="cb6-23"><a href="#cb6-23" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span></span>
<span id="cb6-24"><a href="#cb6-24" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb6-25"><a href="#cb6-25" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>The <code>add</code> instruction demonstrates the typical use of the
R-type format in MIPS instruction set architecture, involving fetching
the instruction, decoding it, executing the operation in the ALU,
skipping memory access, and finally writing back the result to the
register file.</p>
<h3 id="addi-i-type-instruction">ADDI (I-type instruction)</h3>
<p>Name (format, op, function): <code>add immediate (I,8,na)</code></p>
<p>Syntax: <code>addi rt,rs,imm</code></p>
<p>Operation: <code>reg(rt) := reg(rs) + signext(imm);</code></p>
<h4 id="instruction-overview-1">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction from memory.</li>
<li><strong>ID:</strong> Decode the instruction; read the source
register (<code>rs</code>).</li>
<li><strong>EX:</strong> ALU adds the value in the source register to
the immediate value (which is sign-extended).</li>
<li><strong>MEM:</strong> No action.</li>
<li><strong>WB:</strong> The result is written back to the target
register (<code>rt</code>).</li>
</ul>
<h4 id="operation-breakdown-1">Operation Breakdown:</h4>
<p>The following further breaks down the operation of the
<code>ADDI</code> instruction in a MIPS processor across the various
stages of the processor pipeline.</p>
<h5 id="stages-of-the-addi-instruction">Stages of the <code>ADDI</code>
Instruction</h5>
<ol type="1">
<li><strong>Instruction Fetch (IF):</strong>
<ul>
<li>The processor retrieves the <code>ADDI</code> instruction from
memory based on the current Program Counter (PC) value.</li>
<li>The instruction is then forwarded to the next stage for
decoding.</li>
</ul></li>
<li><strong>Instruction Decode (ID):</strong>
<ul>
<li>The instruction is decoded to identify that it is a
<code>ADDI</code> operation.</li>
<li>The source register (<code>rs</code>) is read to obtain its value.
The immediate value (<code>imm</code>) is also extracted from the
instruction during this phase.</li>
</ul></li>
<li><strong>Execute (EX):</strong>
<ul>
<li>The Arithmetic Logic Unit (ALU) performs the addition operation. It
adds the value retrieved from the source register (<code>reg(rs)</code>)
to the sign-extended immediate value (<code>signext(imm)</code>).</li>
<li>This computation involves extending the immediate value to match the
register size (typically 32 bits in MIPS), preserving its sign to handle
negative numbers correctly.</li>
</ul></li>
<li><strong>Memory Access (MEM):</strong>
<ul>
<li>The <code>ADDI</code> instruction does not involve any memory
access, so this stage is effectively a no-op (no operation) for this
instruction.</li>
</ul></li>
<li><strong>Write Back (WB):</strong>
<ul>
<li>The result of the addition from the ALU is written back to the
destination register (<code>reg(rt)</code>).</li>
<li>This step updates the target register with the computed value,
completing the execution of the instruction.</li>
</ul></li>
</ol>
<h5 id="explanation-of-the-code-implementation">Explanation of the Code
Implementation</h5>
<p>The operation of <code>ADDI</code> can be modeled in a simulated or
actual MIPS processor using the following Verilog-like pseudocode:</p>
<div class="sourceCode" id="cb7"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb7-1"><a href="#cb7-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> addi_instruction<span class="op">(</span>rs<span class="op">,</span> rt<span class="op">,</span> imm<span class="op">,</span> <span class="dt">output</span> rt_value<span class="op">);</span></span>
<span id="cb7-2"><a href="#cb7-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs<span class="op">,</span> rt<span class="op">;</span>       <span class="co">// Source and target register indices (5 bits each)</span></span>
<span id="cb7-3"><a href="#cb7-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> imm<span class="op">;</span>         <span class="co">// 16-bit immediate value</span></span>
<span id="cb7-4"><a href="#cb7-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt_value<span class="op">;</span>   <span class="co">// Output to target register</span></span>
<span id="cb7-5"><a href="#cb7-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs_value<span class="op">;</span>     <span class="co">// Value from source register</span></span>
<span id="cb7-6"><a href="#cb7-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> extended_imm<span class="op">;</span> <span class="co">// Sign-extended immediate value</span></span>
<span id="cb7-7"><a href="#cb7-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> extended_imm <span class="op">=</span> <span class="op">{</span></span>
<span id="cb7-8"><a href="#cb7-8" aria-hidden="true" tabindex="-1"></a>      <span class="op">{</span><span class="dv">16</span><span class="op">{</span>imm<span class="op">[</span><span class="dv">15</span><span class="op">]}},</span> imm      <span class="co">// Sign-extend the immediate value</span></span>
<span id="cb7-9"><a href="#cb7-9" aria-hidden="true" tabindex="-1"></a>    <span class="op">};</span></span>
<span id="cb7-10"><a href="#cb7-10" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> rt_value <span class="op">=</span> rs_value <span class="op">+</span> extended_imm<span class="op">;</span></span>
<span id="cb7-11"><a href="#cb7-11" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<ul>
<li><code>rs</code> and <code>rt</code> are inputs representing the
source and destination register indices.</li>
<li><code>imm</code> is the 16-bit immediate value input.</li>
<li>The immediate value is sign-extended to 32 bits using Verilog’s bit
replication and concatenation (<code>{{16{imm[15]}}, imm}</code>), where
<code>imm[15]</code> is the most significant bit (MSB) of the immediate
value, replicated 16 times to fill the upper half of a 32-bit word.</li>
<li>The sum of the sign-extended immediate and the source register value
is computed and assigned to <code>rt_value</code>, which would be
written back to the register file in the actual processor hardware.</li>
</ul>
<h3 id="lw-load-word">LW (Load Word)</h3>
<p>Name (format, op, function): <code>load word (I,35,na)</code></p>
<p>Syntax: <code>lw rt,imm(rs)</code></p>
<p>Operation: <code>reg(rt) := mem[reg(rs) + signext(imm)];</code></p>
<h4 id="instruction-overview-2">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction.</li>
<li><strong>ID:</strong> Decode the instruction; read the base address
register (<code>rs</code>).</li>
<li><strong>EX:</strong> Calculate the memory address by adding the
immediate value (offset) to the base register.</li>
<li><strong>MEM:</strong> Access the memory at the computed address and
read the word.</li>
<li><strong>WB:</strong> Write the loaded word into the target register
(<code>rt</code>).</li>
</ul>
<h4 id="operation-breakdown-2">Operation Breakdown:</h4>
<p>The following provides a detailed breakdown of the operation of the
<code>lw</code> instruction in a MIPS processor, focusing on the key
stages of the processor pipeline.</p>
<h5 id="breakdown-of-lw-instruction-execution">Breakdown of
<code>lw</code> Instruction Execution</h5>
<ol type="1">
<li><strong>Instruction Fetch (IF)</strong>:
<ul>
<li>The processor fetches the <code>lw</code> instruction from
instruction memory using the program counter (PC).</li>
<li>Code snippet showing fetching the instruction from memory:</li>
</ul>
<div class="sourceCode" id="cb8"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb8-1"><a href="#cb8-1" aria-hidden="true" tabindex="-1"></a>  i_Instruction <span class="op">=</span> Imem<span class="op">[</span>i_Addr<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span></code></pre></div></li>
<li><strong>Instruction Decode (ID)</strong>:
<ul>
<li>The fetched instruction is decoded to extract the opcode, source
register (<code>rs</code>), target register (<code>rt</code>), and the
immediate value.</li>
<li>The base address (content of <code>rs</code>) is read from the
register file during this phase.</li>
<li>Code snippet showing the decoding and reading of the base
address:</li>
</ul>
<div class="sourceCode" id="cb9"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb9-1"><a href="#cb9-1" aria-hidden="true" tabindex="-1"></a>  read_data1 <span class="op">=</span> RegData<span class="op">[</span>i_rs<span class="op">];</span>  <span class="co">// Assume i_rs is the source register index</span></span></code></pre></div></li>
<li><strong>Execute (EX)</strong>:
<ul>
<li>The effective memory address is calculated by adding the
sign-extended immediate value to the base address read from
<code>rs</code>.</li>
<li>This calculation typically happens in the ALU.</li>
<li>Code snippet that could represent the address calculation in ALU
(not specifically shown in your snippets):</li>
</ul>
<div class="sourceCode" id="cb10"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb10-1"><a href="#cb10-1" aria-hidden="true" tabindex="-1"></a>  address <span class="op">=</span> read_data1 <span class="op">+</span> sign_extend<span class="op">(</span>imm<span class="op">);</span>  <span class="co">// Conceptual code</span></span></code></pre></div></li>
<li><strong>Memory Access (MEM)</strong>:
<ul>
<li>The processor accesses the memory location computed in the Execute
stage.</li>
<li>The word at this memory address is read.</li>
<li>Code snippet showing memory access to read data:</li>
</ul>
<div class="sourceCode" id="cb11"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb11-1"><a href="#cb11-1" aria-hidden="true" tabindex="-1"></a>  <span class="kw">if</span> <span class="op">(</span>i_MemRead <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="op">{</span></span>
<span id="cb11-2"><a href="#cb11-2" aria-hidden="true" tabindex="-1"></a>    o_rData <span class="op">=</span> Dmem<span class="op">[</span>i_addr<span class="op">];</span>  <span class="co">// Read memory at calculated address</span></span>
<span id="cb11-3"><a href="#cb11-3" aria-hidden="true" tabindex="-1"></a>  <span class="op">}</span></span></code></pre></div></li>
<li><strong>Write Back (WB)</strong>:
<ul>
<li>The data retrieved from memory is written into the target register
(<code>rt</code>).</li>
<li>Code snippet showing the write-back to the register:</li>
</ul>
<div class="sourceCode" id="cb12"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb12-1"><a href="#cb12-1" aria-hidden="true" tabindex="-1"></a>  RegData<span class="op">[</span>i_rt<span class="op">]</span> <span class="op">&lt;=</span> i_wData<span class="op">;</span>  <span class="co">// Assume i_rt is the target register index and i_wData is data read from memory</span></span></code></pre></div></li>
</ol>
<h4 id="sw-store-word">SW (Store Word)</h4>
<p>Name (format, op, function): store word (I,43,na)</p>
<p>Syntax: <code>sw rt,imm(rs)</code></p>
<p>Operation: <code>mem[reg(rs) + signext(imm)] := reg(rt);</code></p>
<h4 id="instruction-overview-3">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction.</li>
<li><strong>ID:</strong> Decode the instruction; read the base address
register (<code>rs</code>) and the register to be stored
(<code>rt</code>).</li>
<li><strong>EX:</strong> Calculate the memory address by adding the
immediate value (offset) to the base register.</li>
<li><strong>MEM:</strong> Write the value from <code>rt</code> into the
calculated memory address.</li>
<li><strong>WB:</strong> No write-back step for store instructions.</li>
</ul>
<h4 id="operation-breakdown-3">Operation Breakdown:</h4>
<p>The <code>SW</code> instruction in the MIPS architecture is used to
store a 32-bit word from a register into memory. Here’s an in-depth
breakdown of how the <code>SW</code> instruction is executed across the
various stages in a MIPS processor.</p>
<h4 id="instruction-stages">Instruction Stages:</h4>
<ol type="1">
<li><strong>IF (Instruction Fetch):</strong>
<ul>
<li>The instruction is fetched from the instruction memory using the
current Program Counter (PC).</li>
</ul></li>
<li><strong>ID (Instruction Decode):</strong>
<ul>
<li>The instruction is decoded to identify it as a <code>SW</code>
instruction.</li>
<li>The base address register (<code>rs</code>) and the register
containing data to be stored (<code>rt</code>) are identified and
read.</li>
</ul></li>
<li><strong>EX (Execute):</strong>
<ul>
<li>The effective memory address is calculated by adding the
sign-extended immediate (offset) to the value in the base register
(<code>rs</code>).</li>
</ul></li>
<li><strong>MEM (Memory Access):</strong>
<ul>
<li>The data in register <code>rt</code> is written to the calculated
memory address.</li>
</ul></li>
<li><strong>WB (Write Back):</strong>
<ul>
<li>No write-back is performed for the <code>SW</code> instruction, as
this instruction does not modify any register contents.</li>
</ul></li>
</ol>
<h5 id="verilog-implementation">Verilog Implementation:</h5>
<p><strong>Data Memory Module (<code>DataMemory.v</code>):</strong></p>
<p>Below is a simplified Verilog module for a data memory component that
can be used to store and retrieve data in a MIPS processor. This module
includes logic for both read and write operations.</p>
<div class="sourceCode" id="cb13"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb13-1"><a href="#cb13-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> DataMemory <span class="op">(</span></span>
<span id="cb13-2"><a href="#cb13-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> clk<span class="op">,</span></span>
<span id="cb13-3"><a href="#cb13-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> memWrite<span class="op">,</span></span>
<span id="cb13-4"><a href="#cb13-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> address<span class="op">,</span></span>
<span id="cb13-5"><a href="#cb13-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> writeData<span class="op">,</span></span>
<span id="cb13-6"><a href="#cb13-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> readData</span>
<span id="cb13-7"><a href="#cb13-7" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb13-8"><a href="#cb13-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> memory <span class="op">[</span><span class="dv">0</span><span class="op">:</span><span class="dv">1023</span><span class="op">];</span></span>
<span id="cb13-9"><a href="#cb13-9" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb13-10"><a href="#cb13-10" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>memWrite<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb13-11"><a href="#cb13-11" aria-hidden="true" tabindex="-1"></a>            memory<span class="op">[</span>address <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">]</span> <span class="op">&lt;=</span> writeData<span class="op">;</span>  <span class="co">// Write operation</span></span>
<span id="cb13-12"><a href="#cb13-12" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb13-13"><a href="#cb13-13" aria-hidden="true" tabindex="-1"></a>            readData <span class="op">&lt;=</span> memory<span class="op">[</span>address <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">];</span>   <span class="co">// Read operation</span></span>
<span id="cb13-14"><a href="#cb13-14" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span></span>
<span id="cb13-15"><a href="#cb13-15" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb13-16"><a href="#cb13-16" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p><strong>Processor Control Logic
(<code>ProcessorControl.v</code>):</strong></p>
<p>Below is an extract from the ProcessorControl module that controls
the behavior of the processor based on the opcode of the instruction
being executed. This snippet shows how the control signals are set for
the <code>SW</code> instruction.</p>
<div class="sourceCode" id="cb14"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb14-1"><a href="#cb14-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ProcessorControl <span class="op">(</span></span>
<span id="cb14-2"><a href="#cb14-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> opcode<span class="op">,</span></span>
<span id="cb14-3"><a href="#cb14-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> memWrite<span class="op">,</span></span>
<span id="cb14-4"><a href="#cb14-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> aluSrc<span class="op">,</span></span>
<span id="cb14-5"><a href="#cb14-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> regDst<span class="op">,</span></span>
<span id="cb14-6"><a href="#cb14-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> memToReg<span class="op">,</span></span>
<span id="cb14-7"><a href="#cb14-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> regWrite</span>
<span id="cb14-8"><a href="#cb14-8" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb14-9"><a href="#cb14-9" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb14-10"><a href="#cb14-10" aria-hidden="true" tabindex="-1"></a>        <span class="kw">case</span> <span class="op">(</span>opcode<span class="op">)</span></span>
<span id="cb14-11"><a href="#cb14-11" aria-hidden="true" tabindex="-1"></a>            <span class="bn">6&#39;b101011</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// Opcode for SW</span></span>
<span id="cb14-12"><a href="#cb14-12" aria-hidden="true" tabindex="-1"></a>                memWrite <span class="op">=</span> <span class="bn">1&#39;b1</span><span class="op">;</span></span>
<span id="cb14-13"><a href="#cb14-13" aria-hidden="true" tabindex="-1"></a>                aluSrc <span class="op">=</span> <span class="bn">1&#39;b1</span><span class="op">;</span></span>
<span id="cb14-14"><a href="#cb14-14" aria-hidden="true" tabindex="-1"></a>                regDst <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-15"><a href="#cb14-15" aria-hidden="true" tabindex="-1"></a>                memToReg <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-16"><a href="#cb14-16" aria-hidden="true" tabindex="-1"></a>                regWrite <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-17"><a href="#cb14-17" aria-hidden="true" tabindex="-1"></a>            <span class="kw">end</span></span>
<span id="cb14-18"><a href="#cb14-18" aria-hidden="true" tabindex="-1"></a>            <span class="kw">default</span><span class="op">:</span> <span class="kw">begin</span></span>
<span id="cb14-19"><a href="#cb14-19" aria-hidden="true" tabindex="-1"></a>                memWrite <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-20"><a href="#cb14-20" aria-hidden="true" tabindex="-1"></a>                aluSrc <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-21"><a href="#cb14-21" aria-hidden="true" tabindex="-1"></a>                regDst <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-22"><a href="#cb14-22" aria-hidden="true" tabindex="-1"></a>                memToReg <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-23"><a href="#cb14-23" aria-hidden="true" tabindex="-1"></a>                regWrite <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb14-24"><a href="#cb14-24" aria-hidden="true" tabindex="-1"></a>            <span class="kw">end</span></span>
<span id="cb14-25"><a href="#cb14-25" aria-hidden="true" tabindex="-1"></a>        <span class="kw">endcase</span></span>
<span id="cb14-26"><a href="#cb14-26" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb14-27"><a href="#cb14-27" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p><strong>Simplified Top-Level MIPS Module:</strong></p>
<div class="sourceCode" id="cb15"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb15-1"><a href="#cb15-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> MIPSProcessor <span class="op">(</span></span>
<span id="cb15-2"><a href="#cb15-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> clk<span class="op">,</span></span>
<span id="cb15-3"><a href="#cb15-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> reset<span class="op">,</span></span>
<span id="cb15-4"><a href="#cb15-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> pc<span class="op">,</span></span>
<span id="cb15-5"><a href="#cb15-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> instruction<span class="op">,</span></span>
<span id="cb15-6"><a href="#cb15-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> aluResult<span class="op">,</span></span>
<span id="cb15-7"><a href="#cb15-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> writeData<span class="op">,</span></span>
<span id="cb15-8"><a href="#cb15-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> readData</span>
<span id="cb15-9"><a href="#cb15-9" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb15-10"><a href="#cb15-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> opcode <span class="op">=</span> instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb15-11"><a href="#cb15-11" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs <span class="op">=</span> instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb15-12"><a href="#cb15-12" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt <span class="op">=</span> instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb15-13"><a href="#cb15-13" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> imm <span class="op">=</span> instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb15-14"><a href="#cb15-14" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> signExtImm <span class="op">=</span> <span class="op">{{</span><span class="dv">16</span><span class="op">{</span>imm<span class="op">[</span><span class="dv">15</span><span class="op">]}},</span> imm<span class="op">};</span></span>
<span id="cb15-15"><a href="#cb15-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> regDataRs<span class="op">,</span> regDataRt<span class="op">;</span></span>
<span id="cb15-16"><a href="#cb15-16" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> memWrite<span class="op">,</span> aluSrc<span class="op">,</span> regDst<span class="op">,</span> memToReg<span class="op">,</span> regWrite<span class="op">;</span></span>
<span id="cb15-17"><a href="#cb15-17" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Instantiate control logic</span></span>
<span id="cb15-18"><a href="#cb15-18" aria-hidden="true" tabindex="-1"></a>    ProcessorControl control<span class="op">(</span>opcode<span class="op">,</span> memWrite<span class="op">,</span> aluSrc<span class="op">,</span> regDst<span class="op">,</span> memToReg<span class="op">,</span> regWrite<span class="op">);</span></span>
<span id="cb15-19"><a href="#cb15-19" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ALU operation (assuming already instantiated and connected)</span></span>
<span id="cb15-20"><a href="#cb15-20" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Data memory operation</span></span>
<span id="cb15-21"><a href="#cb15-21" aria-hidden="true" tabindex="-1"></a>    DataMemory dataMem<span class="op">(</span>clk<span class="op">,</span> memWrite<span class="op">,</span> aluResult<span class="op">,</span> regDataRt<span class="op">,</span> readData<span class="op">);</span></span>
<span id="cb15-22"><a href="#cb15-22" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Register file operations and other connections would be defined here</span></span>
<span id="cb15-23"><a href="#cb15-23" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h3 id="beq-branch-if-equal">BEQ (Branch if Equal)</h3>
<p>Name (format, op, function):
<code>branch on equal (I,4,na)</code></p>
<p>Syntax: <code>beq rs,rt,label</code></p>
<p>Operation: if reg(rs) == reg(rt) then PC = BTA else NOP;</p>
<h4 id="instruction-overview-4">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction.</li>
<li><strong>ID:</strong> Decode the instruction; read the two registers
(<code>rs</code> and <code>rt</code>) and compare them.</li>
<li><strong>EX:</strong> Calculate the branch target address if the
comparison is equal (by adding the sign-extended, shifted immediate to
the PC).</li>
<li><strong>MEM:</strong> No memory access.</li>
<li><strong>WB:</strong> No write-back; update the PC to the branch
address if the condition is met, otherwise increment the PC as
usual.</li>
</ul>
<h4 id="operation-breakdown-4">Operation Breakdown:</h4>
<p>The following provides a detailed breakdown of the operation of the
<code>BEQ</code> instruction in a MIPS processor, focusing on the key
stages of the processor pipeline.</p>
<p>The <code>BEQ</code> (Branch if Equal) instruction in the MIPS
architecture follows a specific flow through the processor stages.
Here’s a step-by-step walkthrough of each stage using Verilog code
examples to illustrate how each part of the instruction’s lifecycle is
handled in hardware.</p>
<h5 id="instruction-stages-for-beq">Instruction Stages for
<code>BEQ</code></h5>
<ol type="1">
<li><strong>IF (Instruction Fetch) Stage</strong>:
<ul>
<li>The instruction is fetched from the instruction memory using the
current Program Counter (PC).</li>
<li>The PC is incremented to point to the next instruction (PC = PC +
4).</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb16"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb16-1"><a href="#cb16-1" aria-hidden="true" tabindex="-1"></a> <span class="kw">module</span> InstructionFetch<span class="op">(</span></span>
<span id="cb16-2"><a href="#cb16-2" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> clk<span class="op">,</span></span>
<span id="cb16-3"><a href="#cb16-3" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> reset<span class="op">,</span></span>
<span id="cb16-4"><a href="#cb16-4" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> next_pc<span class="op">,</span></span>
<span id="cb16-5"><a href="#cb16-5" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> instr<span class="op">,</span></span>
<span id="cb16-6"><a href="#cb16-6" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> pc</span>
<span id="cb16-7"><a href="#cb16-7" aria-hidden="true" tabindex="-1"></a> <span class="op">);</span></span>
<span id="cb16-8"><a href="#cb16-8" aria-hidden="true" tabindex="-1"></a>     <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk <span class="dt">or</span> <span class="kw">posedge</span> reset<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb16-9"><a href="#cb16-9" aria-hidden="true" tabindex="-1"></a>         <span class="kw">if</span> <span class="op">(</span>reset<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb16-10"><a href="#cb16-10" aria-hidden="true" tabindex="-1"></a>             pc <span class="op">&lt;=</span> <span class="bn">32&#39;h00000000</span><span class="op">;</span> <span class="co">// Reset PC to start</span></span>
<span id="cb16-11"><a href="#cb16-11" aria-hidden="true" tabindex="-1"></a>         <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb16-12"><a href="#cb16-12" aria-hidden="true" tabindex="-1"></a>             pc <span class="op">&lt;=</span> next_pc<span class="op">;</span> <span class="co">// Update PC to next PC</span></span>
<span id="cb16-13"><a href="#cb16-13" aria-hidden="true" tabindex="-1"></a>             instr <span class="op">&lt;=</span> instruction_memory<span class="op">[</span>pc <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">];</span> <span class="co">// Fetch instruction from memory</span></span>
<span id="cb16-14"><a href="#cb16-14" aria-hidden="true" tabindex="-1"></a>         <span class="kw">end</span></span>
<span id="cb16-15"><a href="#cb16-15" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span>
<span id="cb16-16"><a href="#cb16-16" aria-hidden="true" tabindex="-1"></a> <span class="kw">endmodule</span></span></code></pre></div>
<ol start="2" type="1">
<li><strong>ID (Instruction Decode) Stage</strong>:
<ul>
<li>Decode the fetched instruction to identify it as
<code>BEQ</code>.</li>
<li>Read the two source registers (<code>rs</code> and <code>rt</code>)
based on the instruction fields.</li>
<li>Set up the control signals for the ALU to perform a subtraction
(<code>rs - rt</code>).</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb17"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb17-1"><a href="#cb17-1" aria-hidden="true" tabindex="-1"></a> <span class="kw">module</span> InstructionDecode<span class="op">(</span></span>
<span id="cb17-2"><a href="#cb17-2" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> instr<span class="op">,</span></span>
<span id="cb17-3"><a href="#cb17-3" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs<span class="op">,</span></span>
<span id="cb17-4"><a href="#cb17-4" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt<span class="op">,</span></span>
<span id="cb17-5"><a href="#cb17-5" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> immediate</span>
<span id="cb17-6"><a href="#cb17-6" aria-hidden="true" tabindex="-1"></a> <span class="op">);</span></span>
<span id="cb17-7"><a href="#cb17-7" aria-hidden="true" tabindex="-1"></a>     <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb17-8"><a href="#cb17-8" aria-hidden="true" tabindex="-1"></a>         rs <span class="op">=</span> instr<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb17-9"><a href="#cb17-9" aria-hidden="true" tabindex="-1"></a>         rt <span class="op">=</span> instr<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb17-10"><a href="#cb17-10" aria-hidden="true" tabindex="-1"></a>         immediate <span class="op">=</span> instr<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span> <span class="co">// For branch offset</span></span>
<span id="cb17-11"><a href="#cb17-11" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span>
<span id="cb17-12"><a href="#cb17-12" aria-hidden="true" tabindex="-1"></a> <span class="kw">endmodule</span></span></code></pre></div>
<ol start="3" type="1">
<li><strong>EX (Execute) Stage</strong>:
<ul>
<li>Compute the target address for branching by sign-extending the
immediate field and shifting left by 2 bits (since it’s word-aligned),
then adding this to the PC + 4 (already incremented PC from IF
stage).</li>
<li>ALU checks if <code>rs</code> and <code>rt</code> are equal by
subtracting and checking if the result is zero.</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb18"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb18-1"><a href="#cb18-1" aria-hidden="true" tabindex="-1"></a> <span class="kw">module</span> ALU<span class="op">(</span></span>
<span id="cb18-2"><a href="#cb18-2" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs_val<span class="op">,</span></span>
<span id="cb18-3"><a href="#cb18-3" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt_val<span class="op">,</span></span>
<span id="cb18-4"><a href="#cb18-4" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> sign_ext_imm<span class="op">,</span></span>
<span id="cb18-5"><a href="#cb18-5" aria-hidden="true" tabindex="-1"></a>     <span class="dt">input</span> <span class="op">[</span><span class="dv">2</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> alu_control<span class="op">,</span></span>
<span id="cb18-6"><a href="#cb18-6" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> zero<span class="op">,</span></span>
<span id="cb18-7"><a href="#cb18-7" aria-hidden="true" tabindex="-1"></a>     <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> alu_result</span>
<span id="cb18-8"><a href="#cb18-8" aria-hidden="true" tabindex="-1"></a> <span class="op">);</span></span>
<span id="cb18-9"><a href="#cb18-9" aria-hidden="true" tabindex="-1"></a>     <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> branch_target <span class="op">=</span> <span class="op">(</span>sign_ext_imm <span class="op">&lt;&lt;</span> <span class="dv">2</span><span class="op">)</span> <span class="op">+</span> pc_plus_4<span class="op">;</span></span>
<span id="cb18-10"><a href="#cb18-10" aria-hidden="true" tabindex="-1"></a>     <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb18-11"><a href="#cb18-11" aria-hidden="true" tabindex="-1"></a>         <span class="kw">case</span><span class="op">(</span>alu_control<span class="op">)</span></span>
<span id="cb18-12"><a href="#cb18-12" aria-hidden="true" tabindex="-1"></a>             <span class="bn">3&#39;b010</span><span class="op">:</span> <span class="kw">begin</span> <span class="co">// Subtract for BEQ</span></span>
<span id="cb18-13"><a href="#cb18-13" aria-hidden="true" tabindex="-1"></a>                 alu_result <span class="op">=</span> rs_val <span class="op">-</span> rt_val<span class="op">;</span></span>
<span id="cb18-14"><a href="#cb18-14" aria-hidden="true" tabindex="-1"></a>                 zero <span class="op">=</span> <span class="op">(</span>alu_result <span class="op">==</span> <span class="dv">0</span><span class="op">)</span> <span class="op">?</span> <span class="bn">1&#39;b1</span> <span class="op">:</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb18-15"><a href="#cb18-15" aria-hidden="true" tabindex="-1"></a>             <span class="kw">end</span></span>
<span id="cb18-16"><a href="#cb18-16" aria-hidden="true" tabindex="-1"></a>         <span class="kw">endcase</span></span>
<span id="cb18-17"><a href="#cb18-17" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span>
<span id="cb18-18"><a href="#cb18-18" aria-hidden="true" tabindex="-1"></a> <span class="kw">endmodule</span></span></code></pre></div>
<ol start="4" type="1">
<li><strong>MEM (Memory Access) Stage</strong>:
<ul>
<li>For <code>BEQ</code>, there is no memory access or data memory
operation.</li>
</ul></li>
<li><strong>WB (Write-Back) Stage</strong>:
<ul>
<li>Update the PC to the branch target if <code>rs</code> ==
<code>rt</code> (if zero flag from ALU is true).</li>
<li>If <code>rs</code> != <code>rt</code>, increment the PC to the next
instruction (already done in IF).</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb19"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb19-1"><a href="#cb19-1" aria-hidden="true" tabindex="-1"></a> <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb19-2"><a href="#cb19-2" aria-hidden="true" tabindex="-1"></a>     <span class="kw">if</span> <span class="op">(</span>branch_taken<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb19-3"><a href="#cb19-3" aria-hidden="true" tabindex="-1"></a>         pc <span class="op">&lt;=</span> branch_target<span class="op">;</span> <span class="co">// Update PC if branch is taken</span></span>
<span id="cb19-4"><a href="#cb19-4" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span>
<span id="cb19-5"><a href="#cb19-5" aria-hidden="true" tabindex="-1"></a> <span class="kw">endmodule</span></span></code></pre></div>
<h5 id="example-verilog-for-complete-beq-control">Example Verilog for
Complete BEQ Control</h5>
<p>Here’s how a simpler control unit might orchestrate these stages just
for the <code>BEQ</code> instruction in a MIPS processor:</p>
<div class="sourceCode" id="cb20"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb20-1"><a href="#cb20-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ControlUnit<span class="op">(</span></span>
<span id="cb20-2"><a href="#cb20-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> opcode<span class="op">,</span></span>
<span id="cb20-3"><a href="#cb20-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> branch<span class="op">,</span></span>
<span id="cb20-4"><a href="#cb20-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> alu_src<span class="op">,</span></span>
<span id="cb20-5"><a href="#cb20-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">2</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> alu_control</span>
<span id="cb20-6"><a href="#cb20-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb20-7"><a href="#cb20-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb20-8"><a href="#cb20-8" aria-hidden="true" tabindex="-1"></a>        <span class="kw">case</span><span class="op">(</span>opcode<span class="op">)</span></span>
<span id="cb20-9"><a href="#cb20-9" aria-hidden="true" tabindex="-1"></a>            <span class="bn">6&#39;b000100</span><span class="op">:</span> <span class="kw">begin</span> <span class="co">// BEQ</span></span>
<span id="cb20-10"><a href="#cb20-10" aria-hidden="true" tabindex="-1"></a>                branch <span class="op">=</span> <span class="bn">1&#39;b1</span><span class="op">;</span></span>
<span id="cb20-11"><a href="#cb20-11" aria-hidden="true" tabindex="-1"></a>                alu_src <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span> <span class="co">// Use rs, rt directly</span></span>
<span id="cb20-12"><a href="#cb20-12" aria-hidden="true" tabindex="-1"></a>                alu_control <span class="op">=</span> <span class="bn">3&#39;b010</span><span class="op">;</span> <span class="co">// Set ALU to subtract</span></span>
<span id="cb20-13"><a href="#cb20-13" aria-hidden="true" tabindex="-1"></a>            <span class="kw">end</span></span>
<span id="cb20-14"><a href="#cb20-14" aria-hidden="true" tabindex="-1"></a>            <span class="kw">default</span><span class="op">:</span> <span class="kw">begin</span></span>
<span id="cb20-15"><a href="#cb20-15" aria-hidden="true" tabindex="-1"></a>                branch <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb20-16"><a href="#cb20-16" aria-hidden="true" tabindex="-1"></a>                alu_src <span class="op">=</span> <span class="bn">1&#39;b0</span><span class="op">;</span></span>
<span id="cb20-17"><a href="#cb20-17" aria-hidden="true" tabindex="-1"></a>                alu_control <span class="op">=</span> <span class="bn">3&#39;b000</span><span class="op">;</span></span>
<span id="cb20-18"><a href="#cb20-18" aria-hidden="true" tabindex="-1"></a>            <span class="kw">end</span></span>
<span id="cb20-19"><a href="#cb20-19" aria-hidden="true" tabindex="-1"></a>        <span class="kw">endcase</span></span>
<span id="cb20-20"><a href="#cb20-20" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb20-21"><a href="#cb20-21" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h3 id="j-jump">J (Jump)</h3>
<p>Name (format, op, function): <code>jump (J,2,na)</code></p>
<p>Syntax: <code>j</code></p>
<p>Operation: <code>PC := JTA;</code></p>
<h4 id="instruction-overview-5">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction.</li>
<li><strong>ID:</strong> Decode the instruction.</li>
<li><strong>EX:</strong> Calculate the jump target address from the
address field of the instruction.</li>
<li><strong>MEM:</strong> No memory access.</li>
<li><strong>WB:</strong> Update the PC to the jump address.</li>
</ul>
<h4 id="operation-breakdown-5">Operation Breakdown:</h4>
<p>The following provides a detailed breakdown of the operation of the
<code>J</code> instruction in a MIPS processor, focusing on the key
stages of the processor pipeline. The J instruction in MIPS is a jump
instruction that allows the program to continue execution from a
specified address. It is used to alter the flow of control
unconditionally.</p>
<h5 id="instruction-format-and-operation">Instruction Format and
Operation:</h5>
<ul>
<li><strong>Name (format, op, function):</strong>
<code>jump (J,2,na)</code></li>
<li><strong>Syntax:</strong> <code>j target</code></li>
<li><strong>Operation:</strong> <code>PC := JTA;</code> where JTA (Jump
Target Address) is calculated from the instruction itself.</li>
</ul>
<h5 id="stages-of-j-instruction-execution">Stages of J Instruction
Execution</h5>
<p>Here’s a breakdown of how the J instruction progresses through each
stage of the MIPS pipeline:</p>
<ol type="1">
<li><strong>IF (Instruction Fetch):</strong>
<ul>
<li>The instruction is fetched from the instruction memory at the
current program counter (PC) address.</li>
<li>The PC is then incremented by 4 to point to the next sequential
instruction (though this increment will be overridden by the jump).</li>
</ul></li>
<li><strong>ID (Instruction Decode):</strong>
<ul>
<li>The opcode of the instruction is decoded to identify it as a jump
instruction.</li>
<li>No registers are read in this stage because the jump instruction
does not involve any registers.</li>
</ul></li>
<li><strong>EX (Execute):</strong>
<ul>
<li>The jump target address (JTA) is calculated from the address field
of the instruction.</li>
<li>JTA is formed by taking the upper 4 bits of the PC (from the
incremented value that points to the next instruction) and concatenating
them with the 26-bit address field from the instruction, shifted left by
2 bits (to word-align the address).</li>
</ul></li>
<li><strong>MEM (Memory Access):</strong>
<ul>
<li>There is no memory access for the jump instruction.</li>
</ul></li>
<li><strong>WB (Write Back):</strong>
<ul>
<li>The PC is updated to the new address calculated in the Execute
stage. This is the jump target address where the program will continue
executing.</li>
</ul></li>
</ol>
<p>Verilog Module for Program Counter with just Jump:</p>
<div class="sourceCode" id="cb21"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb21-1"><a href="#cb21-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ProgramCounter<span class="op">(</span></span>
<span id="cb21-2"><a href="#cb21-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> clk<span class="op">,</span></span>
<span id="cb21-3"><a href="#cb21-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> reset<span class="op">,</span></span>
<span id="cb21-4"><a href="#cb21-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> jump_address<span class="op">,</span>  <span class="co">// Jump target address input</span></span>
<span id="cb21-5"><a href="#cb21-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> jump<span class="op">,</span>                <span class="co">// Control signal to indicate a jump</span></span>
<span id="cb21-6"><a href="#cb21-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> pc       <span class="co">// Program counter output</span></span>
<span id="cb21-7"><a href="#cb21-7" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb21-8"><a href="#cb21-8" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk <span class="dt">or</span> <span class="kw">posedge</span> reset<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb21-9"><a href="#cb21-9" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>reset<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb21-10"><a href="#cb21-10" aria-hidden="true" tabindex="-1"></a>            pc <span class="op">&lt;=</span> <span class="bn">32&#39;b0</span><span class="op">;</span>  <span class="co">// Reset the PC to 0 on reset</span></span>
<span id="cb21-11"><a href="#cb21-11" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span> <span class="kw">else</span> <span class="kw">if</span> <span class="op">(</span>jump<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb21-12"><a href="#cb21-12" aria-hidden="true" tabindex="-1"></a>            pc <span class="op">&lt;=</span> jump_address<span class="op">;</span>  <span class="co">// Update PC to the jump address if jump is asserted</span></span>
<span id="cb21-13"><a href="#cb21-13" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb21-14"><a href="#cb21-14" aria-hidden="true" tabindex="-1"></a>            pc <span class="op">&lt;=</span> pc <span class="op">+</span> <span class="dv">4</span><span class="op">;</span>  <span class="co">// Increment PC by 4 on each clock cycle otherwise</span></span>
<span id="cb21-15"><a href="#cb21-15" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span></span>
<span id="cb21-16"><a href="#cb21-16" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb21-17"><a href="#cb21-17" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h5 id="verilog-module-for-jump-address-calculation">Verilog Module for
Jump Address Calculation:</h5>
<p>While in this project, this function is done by
<code>NextProgramCounter</code> module, here is a simplified version of
a module that calculates the jump address in a MIPS processor to further
illustrate the concept of jump address calculation:</p>
<div class="sourceCode" id="cb22"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb22-1"><a href="#cb22-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> JumpAddressCalculator<span class="op">(</span></span>
<span id="cb22-2"><a href="#cb22-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> address_field<span class="op">,</span>  <span class="co">// Address field from the jump instruction</span></span>
<span id="cb22-3"><a href="#cb22-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> pc_plus_4<span class="op">,</span>      <span class="co">// PC + 4 (the incremented PC pointing to the next instruction)</span></span>
<span id="cb22-4"><a href="#cb22-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> jump_address   <span class="co">// Calculated jump target address</span></span>
<span id="cb22-5"><a href="#cb22-5" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb22-6"><a href="#cb22-6" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> jump_address <span class="op">=</span> <span class="op">{</span>pc_plus_4<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">28</span><span class="op">],</span> address_field <span class="op">&lt;&lt;</span> <span class="dv">2</span><span class="op">};</span></span>
<span id="cb22-7"><a href="#cb22-7" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<ul>
<li>The <code>ProgramCounter</code> module handles updating the PC based
on whether a jump is taken. If a jump is taken, it sets the PC to the
jump address; otherwise, it simply increments the PC.</li>
<li>The <code>JumpAddressCalculator</code> module calculates the full
32-bit jump address by concatenating the upper 4 bits of the incremented
PC (PC+4) with the left-shifted 26-bit address from the jump
instruction.</li>
</ul>
<p>These modules collectively illustrate how the J instruction’s effect
on the program counter can be implemented in hardware using Verilog.</p>
<h3 id="addu-add-unsigned">ADDU (Add Unsigned)</h3>
<p>name (format, op, function): <code>add unsigned (R,0,33)</code></p>
<p>Syntax: <code>addu rd,rs,rt</code></p>
<p>Operation: <code>reg(rd) := reg(rs) + reg(rt);</code></p>
<h4 id="instruction-overview-6">Instruction Overview:</h4>
<ul>
<li><strong>IF (Instruction Fetch):</strong> The instruction is fetched
from memory using the program counter (PC).</li>
<li><strong>ID (Instruction Decode):</strong> The opcode is decoded;
registers rs and rt are read.</li>
<li><strong>EX (Execute):</strong> The arithmetic logic unit (ALU) adds
the values from registers rs and rt.</li>
<li><strong>MEM (Memory Access):</strong> No action needed
(pass-through).</li>
<li><strong>WB (Write Back):</strong> The result from the ALU is written
back to the destination register rd.</li>
</ul>
<h4 id="operation-breakdown-6">Operation Breakdown:</h4>
<p>The following provides a detailed breakdown of the operation of the
<code>ADDU</code> instruction in a MIPS processor, focusing on the key
stages of the processor pipeline.</p>
<p>The <code>ADDU</code> instruction in MIPS is an unsigned addition
operation that does not raise exceptions on overflow. Here’s a detailed
breakdown of how the <code>ADDU</code> instruction is executed across
the various stages of the processor pipeline.</p>
<h5 id="if-instruction-fetch">IF (Instruction Fetch)</h5>
<p>In this stage, the instruction is fetched from the instruction memory
based on the current value of the Program Counter (PC). Here’s how you
might see this operation in Verilog:</p>
<div class="sourceCode" id="cb23"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb23-1"><a href="#cb23-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Instruction Fetch module</span></span>
<span id="cb23-2"><a href="#cb23-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> InstructionFetch<span class="op">(</span></span>
<span id="cb23-3"><a href="#cb23-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_PC<span class="op">,</span>            <span class="co">// Program Counter</span></span>
<span id="cb23-4"><a href="#cb23-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_Instruction   <span class="co">// Fetched instruction</span></span>
<span id="cb23-5"><a href="#cb23-5" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb23-6"><a href="#cb23-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> instruction_memory<span class="op">[</span><span class="dv">255</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span> <span class="co">// Memory array</span></span>
<span id="cb23-7"><a href="#cb23-7" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb23-8"><a href="#cb23-8" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Fetch the instruction</span></span>
<span id="cb23-9"><a href="#cb23-9" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> o_Instruction <span class="op">=</span> instruction_memory<span class="op">[</span>i_PC <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">];</span> <span class="co">// Word aligned access</span></span>
<span id="cb23-10"><a href="#cb23-10" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h5 id="id-instruction-decode">ID (Instruction Decode)</h5>
<p>During this stage, the opcode of the fetched instruction is decoded,
and the register file is accessed to read the contents of registers
<code>rs</code> and <code>rt</code>.</p>
<div class="sourceCode" id="cb24"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb24-1"><a href="#cb24-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Instruction Decode module</span></span>
<span id="cb24-2"><a href="#cb24-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> InstructionDecode<span class="op">(</span></span>
<span id="cb24-3"><a href="#cb24-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Instruction<span class="op">,</span>   <span class="co">// Input from IF stage</span></span>
<span id="cb24-4"><a href="#cb24-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_rs<span class="op">,</span> o_rt<span class="op">,</span> o_rd <span class="co">// Register specifiers</span></span>
<span id="cb24-5"><a href="#cb24-5" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb24-6"><a href="#cb24-6" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Decode the instruction</span></span>
<span id="cb24-7"><a href="#cb24-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> o_rs <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb24-8"><a href="#cb24-8" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> o_rt <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb24-9"><a href="#cb24-9" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> o_rd <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">11</span><span class="op">];</span></span>
<span id="cb24-10"><a href="#cb24-10" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h5 id="ex-execute">EX (Execute)</h5>
<p>The ALU adds the values from registers <code>rs</code> and
<code>rt</code>. Here’s a snippet of the ALU performing this
addition:</p>
<div class="sourceCode" id="cb25"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb25-1"><a href="#cb25-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Arithmetic Logic Unit (ALU) module</span></span>
<span id="cb25-2"><a href="#cb25-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ALU<span class="op">(</span></span>
<span id="cb25-3"><a href="#cb25-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_data1<span class="op">,</span> i_data2<span class="op">,</span>    <span class="co">// Data from registers rs and rt</span></span>
<span id="cb25-4"><a href="#cb25-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">3</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUControl<span class="op">,</span>         <span class="co">// Control signals</span></span>
<span id="cb25-5"><a href="#cb25-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_result        <span class="co">// Result of the ALU operation</span></span>
<span id="cb25-6"><a href="#cb25-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb25-7"><a href="#cb25-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> i_data2<span class="op">,</span> i_ALUControl<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb25-8"><a href="#cb25-8" aria-hidden="true" tabindex="-1"></a>        <span class="kw">case</span> <span class="op">(</span>i_ALUControl<span class="op">)</span></span>
<span id="cb25-9"><a href="#cb25-9" aria-hidden="true" tabindex="-1"></a>            <span class="bn">4&#39;b0010</span><span class="op">:</span> o_result <span class="op">=</span> i_data1 <span class="op">+</span> i_data2<span class="op">;</span> <span class="co">// ADDU operation</span></span>
<span id="cb25-10"><a href="#cb25-10" aria-hidden="true" tabindex="-1"></a>            <span class="co">// Other ALU operations...</span></span>
<span id="cb25-11"><a href="#cb25-11" aria-hidden="true" tabindex="-1"></a>        <span class="kw">endcase</span></span>
<span id="cb25-12"><a href="#cb25-12" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb25-13"><a href="#cb25-13" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h4 id="mem-memory-access">MEM (Memory Access)</h4>
<p>This stage is a pass-through for the ADDU instruction since it does
not involve memory access.</p>
<div class="sourceCode" id="cb26"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb26-1"><a href="#cb26-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Memory Access Stage - No action needed for ADDU</span></span>
<span id="cb26-2"><a href="#cb26-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> MemoryAccess<span class="op">(</span></span>
<span id="cb26-3"><a href="#cb26-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUResult<span class="op">,</span></span>
<span id="cb26-4"><a href="#cb26-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_MemOut</span>
<span id="cb26-5"><a href="#cb26-5" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb26-6"><a href="#cb26-6" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> o_MemOut <span class="op">=</span> i_ALUResult<span class="op">;</span>  <span class="co">// Direct pass-through</span></span>
<span id="cb26-7"><a href="#cb26-7" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h4 id="wb-write-back">WB (Write Back)</h4>
<p>The result from the ALU is written back to the destination register
<code>rd</code>.</p>
<div class="sourceCode" id="cb27"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb27-1"><a href="#cb27-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Write Back stage</span></span>
<span id="cb27-2"><a href="#cb27-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> WriteBack<span class="op">(</span></span>
<span id="cb27-3"><a href="#cb27-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUResult<span class="op">,</span>       <span class="co">// Result from ALU</span></span>
<span id="cb27-4"><a href="#cb27-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_rd<span class="op">,</span>               <span class="co">// Destination register</span></span>
<span id="cb27-5"><a href="#cb27-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_WriteData   <span class="co">// Data to write back</span></span>
<span id="cb27-6"><a href="#cb27-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb27-7"><a href="#cb27-7" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Write the data back to the register file</span></span>
<span id="cb27-8"><a href="#cb27-8" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>i_ALUResult<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb27-9"><a href="#cb27-9" aria-hidden="true" tabindex="-1"></a>        o_WriteData <span class="op">=</span> i_ALUResult<span class="op">;</span></span>
<span id="cb27-10"><a href="#cb27-10" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb27-11"><a href="#cb27-11" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h3 id="sub-subtract">SUB (Subtract)</h3>
<p>Name (format, op, function): <code>subtract (R,0,34)</code></p>
<p>Syntax: <code>sub rd,rs,rt</code></p>
<p>Operation: <code>reg(rd) := reg(rs) [ reg(rt);</code></p>
<h4 id="instruction-overview-7">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction using the PC.</li>
<li><strong>ID:</strong> Decode the instruction; read registers rs and
rt.</li>
<li><strong>EX:</strong> The ALU subtracts the value in rt from rs.</li>
<li><strong>MEM:</strong> No action needed (pass-through).</li>
<li><strong>WB:</strong> The ALU result is written back to register
rd.</li>
</ul>
<h4 id="operation-breakdown-7">Operation Breakdown:</h4>
<p>The following provides a detailed breakdown of the operation of the
<code>SUB</code> instruction in a MIPS processor, focusing on the key
stages of the processor pipeline.</p>
<ol type="1">
<li><strong>Instruction Fetch (IF) Stage:</strong> In this stage, the
processor fetches the instruction from instruction memory using the
Program Counter (PC).</li>
</ol>
<div class="sourceCode" id="cb28"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb28-1"><a href="#cb28-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Instruction Fetch (IF) stage</span></span>
<span id="cb28-2"><a href="#cb28-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> InstructionFetch<span class="op">(</span></span>
<span id="cb28-3"><a href="#cb28-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_pc<span class="op">,</span></span>
<span id="cb28-4"><a href="#cb28-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_instruction</span>
<span id="cb28-5"><a href="#cb28-5" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb28-6"><a href="#cb28-6" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Assume IMem is an array storing instructions</span></span>
<span id="cb28-7"><a href="#cb28-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> IMem<span class="op">[</span><span class="dv">0</span><span class="op">:</span><span class="dv">1023</span><span class="op">];</span></span>
<span id="cb28-8"><a href="#cb28-8" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb28-9"><a href="#cb28-9" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Fetch instruction</span></span>
<span id="cb28-10"><a href="#cb28-10" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>i_pc<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb28-11"><a href="#cb28-11" aria-hidden="true" tabindex="-1"></a>        o_instruction <span class="op">=</span> IMem<span class="op">[</span>i_pc <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">];</span> <span class="co">// Word aligned fetch</span></span>
<span id="cb28-12"><a href="#cb28-12" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb28-13"><a href="#cb28-13" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<ol start="2" type="1">
<li><strong>Instruction Decode (ID) Stage:</strong> Here, the
instruction is decoded, and the relevant registers are read. The
operation is identified, and signals are prepared for the execution
stage.</li>
</ol>
<div class="sourceCode" id="cb29"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb29-1"><a href="#cb29-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Instruction Decode (ID) stage</span></span>
<span id="cb29-2"><a href="#cb29-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> InstructionDecode<span class="op">(</span></span>
<span id="cb29-3"><a href="#cb29-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_instruction<span class="op">,</span></span>
<span id="cb29-4"><a href="#cb29-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_rs<span class="op">,</span> o_rt<span class="op">,</span> o_rd<span class="op">,</span></span>
<span id="cb29-5"><a href="#cb29-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_opcode<span class="op">,</span> o_funct</span>
<span id="cb29-6"><a href="#cb29-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb29-7"><a href="#cb29-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb29-8"><a href="#cb29-8" aria-hidden="true" tabindex="-1"></a>        o_opcode <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb29-9"><a href="#cb29-9" aria-hidden="true" tabindex="-1"></a>        o_rs <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb29-10"><a href="#cb29-10" aria-hidden="true" tabindex="-1"></a>        o_rt <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb29-11"><a href="#cb29-11" aria-hidden="true" tabindex="-1"></a>        o_rd <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">11</span><span class="op">];</span></span>
<span id="cb29-12"><a href="#cb29-12" aria-hidden="true" tabindex="-1"></a>        o_funct <span class="op">=</span> i_instruction<span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb29-13"><a href="#cb29-13" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb29-14"><a href="#cb29-14" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<ol start="3" type="1">
<li><strong>Execution (EX) Stage:</strong> The ALU performs the
subtraction based on the decoded instruction. The operands are taken
from the registers identified in the ID stage.</li>
</ol>
<div class="sourceCode" id="cb30"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb30-1"><a href="#cb30-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Execution (EX) stage - ALU for SUB operation</span></span>
<span id="cb30-2"><a href="#cb30-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ALU<span class="op">(</span></span>
<span id="cb30-3"><a href="#cb30-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_data1<span class="op">,</span> i_data2<span class="op">,</span></span>
<span id="cb30-4"><a href="#cb30-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">3</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUcontrol<span class="op">,</span></span>
<span id="cb30-5"><a href="#cb30-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_result</span>
<span id="cb30-6"><a href="#cb30-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb30-7"><a href="#cb30-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb30-8"><a href="#cb30-8" aria-hidden="true" tabindex="-1"></a>        <span class="kw">case</span><span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb30-9"><a href="#cb30-9" aria-hidden="true" tabindex="-1"></a>            <span class="bn">4&#39;b0110</span><span class="op">:</span> o_result <span class="op">=</span> i_data1 <span class="op">-</span> i_data2<span class="op">;</span> <span class="co">// SUB operation</span></span>
<span id="cb30-10"><a href="#cb30-10" aria-hidden="true" tabindex="-1"></a>            <span class="co">// Additional cases for other ALU operations</span></span>
<span id="cb30-11"><a href="#cb30-11" aria-hidden="true" tabindex="-1"></a>        <span class="kw">endcase</span></span>
<span id="cb30-12"><a href="#cb30-12" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb30-13"><a href="#cb30-13" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<ol start="4" type="1">
<li><strong>Memory (MEM) Stage:</strong> For the <code>SUB</code>
instruction, there is no memory operation needed. This stage can be
passed through or handled with a control signal that disables memory
operations.</li>
</ol>
<div class="sourceCode" id="cb31"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb31-1"><a href="#cb31-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Memory (MEM) stage pass-through for SUB</span></span>
<span id="cb31-2"><a href="#cb31-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> MemoryStage<span class="op">(</span></span>
<span id="cb31-3"><a href="#cb31-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemRead<span class="op">,</span> i_MemWrite<span class="op">,</span></span>
<span id="cb31-4"><a href="#cb31-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_address<span class="op">,</span> i_writeData<span class="op">,</span></span>
<span id="cb31-5"><a href="#cb31-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_readData</span>
<span id="cb31-6"><a href="#cb31-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb31-7"><a href="#cb31-7" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Memory array</span></span>
<span id="cb31-8"><a href="#cb31-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> DMem<span class="op">[</span><span class="dv">0</span><span class="op">:</span><span class="dv">1023</span><span class="op">];</span></span>
<span id="cb31-9"><a href="#cb31-9" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb31-10"><a href="#cb31-10" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb31-11"><a href="#cb31-11" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>i_MemWrite<span class="op">)</span> DMem<span class="op">[</span>i_address <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">]</span> <span class="op">=</span> i_writeData<span class="op">;</span></span>
<span id="cb31-12"><a href="#cb31-12" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>i_MemRead<span class="op">)</span> o_readData <span class="op">=</span> DMem<span class="op">[</span>i_address <span class="op">&gt;&gt;</span> <span class="dv">2</span><span class="op">];</span></span>
<span id="cb31-13"><a href="#cb31-13" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb31-14"><a href="#cb31-14" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<ol start="5" type="1">
<li><strong>Write Back (WB) Stage:</strong> The result of the ALU
operation is written back to the register file, particularly in the
register specified by <code>rd</code>.</li>
</ol>
<div class="sourceCode" id="cb32"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb32-1"><a href="#cb32-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Write Back (WB) Stage</span></span>
<span id="cb32-2"><a href="#cb32-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> WriteBack<span class="op">(</span></span>
<span id="cb32-3"><a href="#cb32-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUresult<span class="op">,</span></span>
<span id="cb32-4"><a href="#cb32-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_rd<span class="op">,</span></span>
<span id="cb32-5"><a href="#cb32-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_RegWrite<span class="op">,</span></span>
<span id="cb32-6"><a href="#cb32-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_writeData</span>
<span id="cb32-7"><a href="#cb32-7" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb32-8"><a href="#cb32-8" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>i_ALUresult<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb32-9"><a href="#cb32-9" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>i_RegWrite<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb32-10"><a href="#cb32-10" aria-hidden="true" tabindex="-1"></a>            o_writeData <span class="op">=</span> i_ALUresult<span class="op">;</span></span>
<span id="cb32-11"><a href="#cb32-11" aria-hidden="true" tabindex="-1"></a>        <span class="kw">end</span></span>
<span id="cb32-12"><a href="#cb32-12" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb32-13"><a href="#cb32-13" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<h3 id="subu-subtract-unsigned">SUBU (Subtract Unsigned)</h3>
<p>Name (format, op, function):
<code>subtract unsigned (R,0,35)</code></p>
<p>Syntax: <code>subu rd,rs,rt</code></p>
<p>Operation: <code>reg(rd) := reg(rs) [ reg(rt);</code></p>
<h4 id="instruction-overview-8">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Fetch the instruction using the PC.</li>
<li><strong>ID:</strong> Decode the instruction; read registers rs and
rt.</li>
<li><strong>EX:</strong> The ALU subtracts the value in rt from rs.</li>
<li><strong>MEM:</strong> No action needed (pass-through).</li>
<li><strong>WB:</strong> The ALU result is written back to register
rd.</li>
</ul>
<h4 id="operation-breakdown-8">Operation Breakdown:</h4>
<ol type="1">
<li><strong>Instruction Fetch (IF)</strong>
<ul>
<li>The instruction is fetched from the instruction memory using the
Program Counter (PC).</li>
<li>Verilog snippet:</li>
</ul>
<div class="sourceCode" id="cb33"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb33-1"><a href="#cb33-1" aria-hidden="true" tabindex="-1"></a>  <span class="co">// IF Stage</span></span>
<span id="cb33-2"><a href="#cb33-2" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb33-3"><a href="#cb33-3" aria-hidden="true" tabindex="-1"></a>      <span class="kw">if</span> <span class="op">(</span>reset<span class="op">)</span></span>
<span id="cb33-4"><a href="#cb33-4" aria-hidden="true" tabindex="-1"></a>          pc <span class="op">&lt;=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb33-5"><a href="#cb33-5" aria-hidden="true" tabindex="-1"></a>      <span class="kw">else</span> <span class="kw">if</span> <span class="op">(</span>pc_src<span class="op">)</span></span>
<span id="cb33-6"><a href="#cb33-6" aria-hidden="true" tabindex="-1"></a>          pc <span class="op">&lt;=</span> pc_next<span class="op">;</span></span>
<span id="cb33-7"><a href="#cb33-7" aria-hidden="true" tabindex="-1"></a>      <span class="kw">else</span></span>
<span id="cb33-8"><a href="#cb33-8" aria-hidden="true" tabindex="-1"></a>          pc <span class="op">&lt;=</span> pc <span class="op">+</span> <span class="dv">4</span><span class="op">;</span></span>
<span id="cb33-9"><a href="#cb33-9" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span></code></pre></div></li>
<li><strong>Instruction Decode (ID)</strong>
<ul>
<li>The fetched instruction is decoded to identify the operation as
<code>SUBU</code> and the source (<code>rs</code>, <code>rt</code>) and
destination (<code>rd</code>) registers are identified.</li>
<li>Verilog snippet:</li>
</ul>
<div class="sourceCode" id="cb34"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb34-1"><a href="#cb34-1" aria-hidden="true" tabindex="-1"></a>  <span class="co">// ID Stage</span></span>
<span id="cb34-2"><a href="#cb34-2" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> instruction<span class="op">;</span></span>
<span id="cb34-3"><a href="#cb34-3" aria-hidden="true" tabindex="-1"></a>  <span class="dt">wire</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs<span class="op">,</span> rt<span class="op">,</span> rd<span class="op">;</span></span>
<span id="cb34-4"><a href="#cb34-4" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> rs <span class="op">=</span> instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb34-5"><a href="#cb34-5" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> rt <span class="op">=</span> instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb34-6"><a href="#cb34-6" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> rd <span class="op">=</span> instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">11</span><span class="op">];</span></span></code></pre></div></li>
<li><strong>Execution (EX)</strong>
<ul>
<li>The actual subtraction of the contents of the registers
<code>rs</code> and <code>rt</code> is performed. The result does not
account for overflow because it is unsigned.</li>
<li>Verilog snippet:</li>
</ul>
<div class="sourceCode" id="cb35"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb35-1"><a href="#cb35-1" aria-hidden="true" tabindex="-1"></a>  <span class="co">// EX Stage</span></span>
<span id="cb35-2"><a href="#cb35-2" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> reg_data<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span>  <span class="co">// Register file</span></span>
<span id="cb35-3"><a href="#cb35-3" aria-hidden="true" tabindex="-1"></a>  <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs_value<span class="op">,</span> rt_value<span class="op">,</span> result<span class="op">;</span></span>
<span id="cb35-4"><a href="#cb35-4" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> rs_value <span class="op">=</span> reg_data<span class="op">[</span>rs<span class="op">];</span></span>
<span id="cb35-5"><a href="#cb35-5" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> rt_value <span class="op">=</span> reg_data<span class="op">[</span>rt<span class="op">];</span></span>
<span id="cb35-6"><a href="#cb35-6" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> result <span class="op">=</span> rs_value <span class="op">-</span> rt_value<span class="op">;</span></span></code></pre></div></li>
<li><strong>Memory Access (MEM)</strong>
<ul>
<li><code>SUBU</code> does not require a memory operation, so this stage
can be considered a pass-through.</li>
<li>No action required for SUBU</li>
</ul></li>
<li><strong>Write Back (WB)</strong>
<ul>
<li>The result of the subtraction is written back to the destination
register <code>rd</code>.</li>
<li>Verilog snippet:</li>
</ul>
<div class="sourceCode" id="cb36"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb36-1"><a href="#cb36-1" aria-hidden="true" tabindex="-1"></a>  <span class="co">// WB Stage</span></span>
<span id="cb36-2"><a href="#cb36-2" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb36-3"><a href="#cb36-3" aria-hidden="true" tabindex="-1"></a>      <span class="kw">if</span> <span class="op">(</span>reg_write<span class="op">)</span></span>
<span id="cb36-4"><a href="#cb36-4" aria-hidden="true" tabindex="-1"></a>          reg_data<span class="op">[</span>rd<span class="op">]</span> <span class="op">&lt;=</span> result<span class="op">;</span></span>
<span id="cb36-5"><a href="#cb36-5" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span></code></pre></div></li>
</ol>
<div class="sourceCode" id="cb37"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb37-1"><a href="#cb37-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> MIPS_Processor<span class="op">(</span><span class="dt">input</span> clk<span class="op">,</span> <span class="dt">input</span> reset<span class="op">,</span> <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> pc<span class="op">);</span></span>
<span id="cb37-2"><a href="#cb37-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> pc<span class="op">,</span> next_pc<span class="op">;</span></span>
<span id="cb37-3"><a href="#cb37-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> reg_file<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span>  <span class="co">// Register file</span></span>
<span id="cb37-4"><a href="#cb37-4" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb37-5"><a href="#cb37-5" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Instruction Fetch</span></span>
<span id="cb37-6"><a href="#cb37-6" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb37-7"><a href="#cb37-7" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>reset<span class="op">)</span></span>
<span id="cb37-8"><a href="#cb37-8" aria-hidden="true" tabindex="-1"></a>            pc <span class="op">&lt;=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb37-9"><a href="#cb37-9" aria-hidden="true" tabindex="-1"></a>        <span class="kw">else</span></span>
<span id="cb37-10"><a href="#cb37-10" aria-hidden="true" tabindex="-1"></a>            pc <span class="op">&lt;=</span> next_pc<span class="op">;</span></span>
<span id="cb37-11"><a href="#cb37-11" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb37-12"><a href="#cb37-12" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb37-13"><a href="#cb37-13" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Instruction Decode</span></span>
<span id="cb37-14"><a href="#cb37-14" aria-hidden="true" tabindex="-1"></a>    <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> instruction<span class="op">;</span></span>
<span id="cb37-15"><a href="#cb37-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs<span class="op">,</span> rt<span class="op">,</span> rd<span class="op">;</span></span>
<span id="cb37-16"><a href="#cb37-16" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> opcode<span class="op">;</span></span>
<span id="cb37-17"><a href="#cb37-17" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> opcode <span class="op">=</span> instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb37-18"><a href="#cb37-18" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> rs <span class="op">=</span> instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb37-19"><a href="#cb37-19" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> rt <span class="op">=</span> instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb37-20"><a href="#cb37-20" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> rd <span class="op">=</span> instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">11</span><span class="op">];</span></span>
<span id="cb37-21"><a href="#cb37-21" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Execute</span></span>
<span id="cb37-22"><a href="#cb37-22" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs_value<span class="op">,</span> rt_value<span class="op">,</span> alu_result<span class="op">;</span></span>
<span id="cb37-23"><a href="#cb37-23" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> rs_value <span class="op">=</span> reg_file<span class="op">[</span>rs<span class="op">];</span></span>
<span id="cb37-24"><a href="#cb37-24" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> rt_value <span class="op">=</span> reg_file<span class="op">[</span>rt<span class="op">];</span></span>
<span id="cb37-25"><a href="#cb37-25" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> alu_result <span class="op">=</span> <span class="op">(</span>opcode <span class="op">==</span> <span class="bn">6&#39;b000000</span><span class="op">)</span> <span class="op">?</span> <span class="op">(</span>rs_value <span class="op">-</span> rt_value<span class="op">)</span> <span class="op">:</span> <span class="bn">32&#39;b0</span><span class="op">;</span>  <span class="co">// SUBU Opcode assumed</span></span>
<span id="cb37-26"><a href="#cb37-26" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Memory Access</span></span>
<span id="cb37-27"><a href="#cb37-27" aria-hidden="true" tabindex="-1"></a>    <span class="co">// No memory access for SUBU</span></span>
<span id="cb37-28"><a href="#cb37-28" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Write Back</span></span>
<span id="cb37-29"><a href="#cb37-29" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb37-30"><a href="#cb37-30" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(</span>opcode <span class="op">==</span> <span class="bn">6&#39;b100011</span><span class="op">)</span>  <span class="co">// SUBU Opcode</span></span>
<span id="cb37-31"><a href="#cb37-31" aria-hidden="true" tabindex="-1"></a>            reg_file<span class="op">[</span>rd<span class="op">]</span> <span class="op">&lt;=</span> alu_result<span class="op">;</span></span>
<span id="cb37-32"><a href="#cb37-32" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb37-33"><a href="#cb37-33" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Program Counter Update</span></span>
<span id="cb37-34"><a href="#cb37-34" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb37-35"><a href="#cb37-35" aria-hidden="true" tabindex="-1"></a>        next_pc <span class="op">=</span> pc <span class="op">+</span> <span class="dv">4</span><span class="op">;</span>  <span class="co">// Simple sequential execution</span></span>
<span id="cb37-36"><a href="#cb37-36" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb37-37"><a href="#cb37-37" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>In this example, <code>opcode == 6'b100011'</code> is the actual
opcode for <code>SUBU</code>.</p>
<h3 id="and-bitwise-and">AND (Bitwise AND)</h3>
<ul>
<li><strong>IF:</strong> Instruction is fetched.</li>
<li><strong>ID:</strong> Instruction is decoded. For AND, rs and rt are
read;</li>
<li><strong>EX:</strong> The ALU performs an AND operation between
operands.</li>
<li><strong>MEM:</strong> No action needed.</li>
<li><strong>WB:</strong> Result is written back to rd (AND)</li>
</ul>
<h4 id="instruction-breakdown">Instruction Breakdown</h4>
<p>The AND instruction performs a bitwise AND operation between two
operands and stores the result in a destination register.</p>
<ol type="1">
<li><strong>Instruction Fetch (IF) stage:</strong>
<ul>
<li>The Program Counter (PC) contains the address of the AND
instruction.</li>
<li>The Instruction Memory module (<code>InstructionMemory.v</code>)
fetches the instruction from the memory location pointed to by the
PC.</li>
<li>The fetched instruction is passed to the next stage.</li>
</ul></li>
<li><strong>Instruction Decode (ID) stage:</strong>
<ul>
<li>The Control Unit module (<code>ControlUnit.v</code>) decodes the
opcode of the AND instruction.</li>
<li>Based on the opcode, the Control Unit generates the appropriate
control signals for the data-path components.</li>
<li>The Register File module reads the values of the source registers
specified in the AND instruction.</li>
</ul></li>
<li><strong>Execution (EX) stage:</strong>
<ul>
<li>The ALU module (<code>ALU.v</code>) performs the bitwise AND
operation between the values of the source registers.</li>
<li>The ALU control signal generated by the Control Unit determines the
specific operation to be performed (AND in this case).</li>
</ul></li>
<li><strong>Memory Access (MEM) stage:</strong>
<ul>
<li>The AND instruction does not require any memory access, so this
stage is a pass-through.</li>
</ul></li>
<li><strong>Write Back (WB) stage:</strong>
<ul>
<li>The result of the AND operation from the ALU is written back to the
destination register specified in the AND instruction.</li>
<li>The RegWrite control signal generated by the Control Unit enables
the writing of the result to the Register File.</li>
</ul></li>
</ol>
<p>Within the control unit, the AND instruction is identified by its
opcode, and the appropriate control signals are set to execute the AND
operation.</p>
<p>The ALU module performs the bitwise AND operation between the source
register values, and the result is written back to the destination
register.</p>
<p>Here’s an example of how the AND instruction flows through the
different stages of the single-cycle MIPS processor starting within the
control unit.</p>
<div class="sourceCode" id="cb38"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb38-1"><a href="#cb38-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Control Unit</span></span>
<span id="cb38-2"><a href="#cb38-2" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb38-3"><a href="#cb38-3" aria-hidden="true" tabindex="-1"></a>  <span class="kw">case</span> <span class="op">(</span>i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">])</span></span>
<span id="cb38-4"><a href="#cb38-4" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb38-5"><a href="#cb38-5" aria-hidden="true" tabindex="-1"></a>    <span class="bn">6&#39;b001100</span><span class="op">:</span> <span class="kw">begin</span>   <span class="co">// andi</span></span>
<span id="cb38-6"><a href="#cb38-6" aria-hidden="true" tabindex="-1"></a>      o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>    <span class="co">// Destination register is rt</span></span>
<span id="cb38-7"><a href="#cb38-7" aria-hidden="true" tabindex="-1"></a>      o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span>    <span class="co">// Second operand is immediate value</span></span>
<span id="cb38-8"><a href="#cb38-8" aria-hidden="true" tabindex="-1"></a>      o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// ALU result is written to register</span></span>
<span id="cb38-9"><a href="#cb38-9" aria-hidden="true" tabindex="-1"></a>      o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span>  <span class="co">// Write to register file</span></span>
<span id="cb38-10"><a href="#cb38-10" aria-hidden="true" tabindex="-1"></a>      o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>   <span class="co">// No memory read</span></span>
<span id="cb38-11"><a href="#cb38-11" aria-hidden="true" tabindex="-1"></a>      o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// No memory write</span></span>
<span id="cb38-12"><a href="#cb38-12" aria-hidden="true" tabindex="-1"></a>      o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>    <span class="co">// No branch</span></span>
<span id="cb38-13"><a href="#cb38-13" aria-hidden="true" tabindex="-1"></a>      o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>       <span class="co">// No branch if not equal</span></span>
<span id="cb38-14"><a href="#cb38-14" aria-hidden="true" tabindex="-1"></a>      o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b11</span><span class="op">;</span> <span class="co">// ALU operation is AND</span></span>
<span id="cb38-15"><a href="#cb38-15" aria-hidden="true" tabindex="-1"></a>      o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>      <span class="co">// No jump</span></span>
<span id="cb38-16"><a href="#cb38-16" aria-hidden="true" tabindex="-1"></a>      <span class="co">// ...</span></span>
<span id="cb38-17"><a href="#cb38-17" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb38-18"><a href="#cb38-18" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb38-19"><a href="#cb38-19" aria-hidden="true" tabindex="-1"></a>  <span class="kw">endcase</span></span>
<span id="cb38-20"><a href="#cb38-20" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span>
<span id="cb38-21"><a href="#cb38-21" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb38-22"><a href="#cb38-22" aria-hidden="true" tabindex="-1"></a><span class="co">// ALU</span></span>
<span id="cb38-23"><a href="#cb38-23" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> data2<span class="op">,</span> i_ALUcontrol<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb38-24"><a href="#cb38-24" aria-hidden="true" tabindex="-1"></a>  <span class="kw">case</span> <span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb38-25"><a href="#cb38-25" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb38-26"><a href="#cb38-26" aria-hidden="true" tabindex="-1"></a>    <span class="bn">4&#39;b0000</span><span class="op">:</span>  <span class="co">// AND</span></span>
<span id="cb38-27"><a href="#cb38-27" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">&amp;</span> data2<span class="op">;</span></span>
<span id="cb38-28"><a href="#cb38-28" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb38-29"><a href="#cb38-29" aria-hidden="true" tabindex="-1"></a>  <span class="kw">endcase</span></span>
<span id="cb38-30"><a href="#cb38-30" aria-hidden="true" tabindex="-1"></a>  <span class="co">// ...</span></span>
<span id="cb38-31"><a href="#cb38-31" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<p>In the Control Unit module, when the opcode of the instruction
matches the AND opcode (<code>6'b001100</code> in this case), the
appropriate control signals are set.</p>
<p>The <code>ALUSrc</code> signal is set to 1 to select the immediate
value as the second operand, and the <code>ALUOp</code> signal is set to
indicate an AND operation.</p>
<p>In the ALU module, when the <code>ALUcontrol</code> signal matches
the AND operation (4’b0000), the bitwise AND operation is performed
between the two input operands (i_data1 and data2), and the result is
assigned to <code>o_ALUresult</code>.</p>
<h3 id="andi-and-immediate">ANDI (AND Immediate)</h3>
<p><code>and immediate (I,12,na)</code></p>
<p><code>andi rt,rs,imm</code></p>
<p><code>reg(rt) := reg(rs) &amp; zeroext(imm);</code></p>
<h4 id="instruction-overview-9">Instruction Overview:</h4>
<ul>
<li><strong>IF:</strong> Instruction is fetched.</li>
<li><strong>ID:</strong> Opcode decoded. Registers rs and immediate for
ANDI are read.</li>
<li><strong>EX:</strong> ALU performs an AND operation.</li>
<li><strong>MEM:</strong> No memory access.</li>
<li><strong>WB:</strong> Result written to or rt (ANDI).</li>
</ul>
<h4 id="instruction-breakdown-1">Instruction Breakdown</h4>
<p>The ANDI (AND Immediate) instruction performs a bitwise AND operation
between a register value and an immediate value.</p>
<p>Here’s an explanation of how the ANDI instruction goes through each
stage of the MIPS processor pipeline:</p>
<ol type="1">
<li><strong>Instruction Fetch (IF):</strong>
<ul>
<li>The Program Counter (PC) contains the address of the ANDI
instruction in the Instruction Memory.</li>
<li>The instruction is fetched from the Instruction Memory using the PC
value.</li>
<li>Example code in the Instruction Memory module
(<code>InstructionMemory.v</code>):</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb39"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb39-1"><a href="#cb39-1" aria-hidden="true" tabindex="-1"></a>     <span class="kw">always</span> <span class="op">@(</span>i_Addr<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb39-2"><a href="#cb39-2" aria-hidden="true" tabindex="-1"></a>       <span class="kw">if</span> <span class="op">(</span>i_Addr <span class="op">==</span> <span class="op">-</span><span class="dv">4</span><span class="op">)</span> <span class="kw">begin</span>         <span class="co">// init</span></span>
<span id="cb39-3"><a href="#cb39-3" aria-hidden="true" tabindex="-1"></a>         i_Instruction <span class="op">=</span> <span class="bn">32&#39;b11111100000000000000000000000000</span><span class="op">;</span></span>
<span id="cb39-4"><a href="#cb39-4" aria-hidden="true" tabindex="-1"></a>       <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb39-5"><a href="#cb39-5" aria-hidden="true" tabindex="-1"></a>         i_Instruction <span class="op">=</span> Imem<span class="op">[</span>i_Addr<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span>
<span id="cb39-6"><a href="#cb39-6" aria-hidden="true" tabindex="-1"></a>       <span class="kw">end</span></span>
<span id="cb39-7"><a href="#cb39-7" aria-hidden="true" tabindex="-1"></a>       i_Ctr <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb39-8"><a href="#cb39-8" aria-hidden="true" tabindex="-1"></a>       i_Funcode <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb39-9"><a href="#cb39-9" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span></code></pre></div>
<ol start="2" type="1">
<li><strong>Instruction Decode (ID):</strong>
<ul>
<li>The fetched instruction is decoded to determine the operation to be
performed.</li>
<li>The Control Unit generates the necessary control signals based on
the opcode and function code of the instruction.</li>
<li>The register to be read (rs) is determined from the instruction, and
the immediate value is sign-extended.</li>
<li>Example code in the Control Unit module
(<code>ControlUnit.v</code>):</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb40"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb40-1"><a href="#cb40-1" aria-hidden="true" tabindex="-1"></a>     <span class="bn">6&#39;b001100</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// andi</span></span>
<span id="cb40-2"><a href="#cb40-2" aria-hidden="true" tabindex="-1"></a>       o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-3"><a href="#cb40-3" aria-hidden="true" tabindex="-1"></a>       o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb40-4"><a href="#cb40-4" aria-hidden="true" tabindex="-1"></a>       o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-5"><a href="#cb40-5" aria-hidden="true" tabindex="-1"></a>       o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb40-6"><a href="#cb40-6" aria-hidden="true" tabindex="-1"></a>       o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-7"><a href="#cb40-7" aria-hidden="true" tabindex="-1"></a>       o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-8"><a href="#cb40-8" aria-hidden="true" tabindex="-1"></a>       o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-9"><a href="#cb40-9" aria-hidden="true" tabindex="-1"></a>       o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-10"><a href="#cb40-10" aria-hidden="true" tabindex="-1"></a>       o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b11</span><span class="op">;</span></span>
<span id="cb40-11"><a href="#cb40-11" aria-hidden="true" tabindex="-1"></a>       o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb40-12"><a href="#cb40-12" aria-hidden="true" tabindex="-1"></a>       <span class="co">// ...</span></span>
<span id="cb40-13"><a href="#cb40-13" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span></code></pre></div>
<ol start="3" type="1">
<li><strong>Execute (EX):</strong>
<ul>
<li>The ALU performs the bitwise AND operation between the value of
register rs and the sign-extended immediate value.</li>
<li>The result of the AND operation is stored in a temporary
register.</li>
<li>Example code in the ALU module (<code>ALU.v</code>):</li>
</ul>
<div class="sourceCode" id="cb41"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb41-1"><a href="#cb41-1" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> data2<span class="op">,</span> i_ALUcontrol<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb41-2"><a href="#cb41-2" aria-hidden="true" tabindex="-1"></a>    <span class="kw">case</span> <span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb41-3"><a href="#cb41-3" aria-hidden="true" tabindex="-1"></a>      <span class="co">// ...</span></span>
<span id="cb41-4"><a href="#cb41-4" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0000</span><span class="op">:</span>  <span class="co">// AND</span></span>
<span id="cb41-5"><a href="#cb41-5" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> i_data1 <span class="op">&amp;</span> data2<span class="op">;</span> <span class="co">// bitwise AND</span></span>
<span id="cb41-6"><a href="#cb41-6" aria-hidden="true" tabindex="-1"></a>      <span class="co">// ...</span></span>
<span id="cb41-7"><a href="#cb41-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">endcase</span></span>
<span id="cb41-8"><a href="#cb41-8" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb41-9"><a href="#cb41-9" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span></code></pre></div></li>
<li><strong>Memory Access (MEM):</strong>
<ul>
<li>The ANDI instruction does not involve memory access, so no action is
needed in this stage.</li>
<li>The result from the Execute stage is simply passed through to the
next stage.</li>
</ul></li>
<li><strong>Write Back (WB):</strong>
<ul>
<li>The result of the AND operation, stored in the temporary register,
is written back to the destination register (rt) in the Register
File.</li>
<li>Example code in the Register File module
(<code>RegisterFile.v</code>):</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb42"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb42-1"><a href="#cb42-1" aria-hidden="true" tabindex="-1"></a>     <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_Clk <span class="dt">or</span> <span class="kw">posedge</span> i_Rst<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb42-2"><a href="#cb42-2" aria-hidden="true" tabindex="-1"></a>       <span class="kw">if</span> <span class="op">(</span>i_Rst<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb42-3"><a href="#cb42-3" aria-hidden="true" tabindex="-1"></a>         <span class="co">// Reset all registers to zero</span></span>
<span id="cb42-4"><a href="#cb42-4" aria-hidden="true" tabindex="-1"></a>         <span class="kw">for</span> <span class="op">(</span>j <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> j <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> j <span class="op">=</span> j <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb42-5"><a href="#cb42-5" aria-hidden="true" tabindex="-1"></a>           RegData<span class="op">[</span>j<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb42-6"><a href="#cb42-6" aria-hidden="true" tabindex="-1"></a>         <span class="kw">end</span></span>
<span id="cb42-7"><a href="#cb42-7" aria-hidden="true" tabindex="-1"></a>       <span class="kw">end</span> <span class="kw">else</span> <span class="kw">if</span> <span class="op">(</span>i_RegWrite<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb42-8"><a href="#cb42-8" aria-hidden="true" tabindex="-1"></a>         <span class="co">// Write data to the specified register</span></span>
<span id="cb42-9"><a href="#cb42-9" aria-hidden="true" tabindex="-1"></a>         RegData<span class="op">[</span>i_wReg<span class="op">]</span> <span class="op">=</span> i_wData<span class="op">;</span></span>
<span id="cb42-10"><a href="#cb42-10" aria-hidden="true" tabindex="-1"></a>       <span class="kw">end</span></span>
<span id="cb42-11"><a href="#cb42-11" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span></code></pre></div>
<p>Here’s an example of how the ANDI instruction would look in machine
code:</p>
<pre><code>001100 01000 01010 0000000000001111</code></pre>
<p>In this example: - <code>001100</code> is the opcode for the ANDI
instruction. - <code>01000</code> represents the source register (rs),
which is $8 in this case. - <code>01010</code> represents the
destination register (rt), which is $10 in this case. -
<code>0000000000001111</code> is the immediate value, which is 15 in
decimal.</p>
<h3 id="or">OR</h3>
<ul>
<li><strong>IF:</strong> Instruction is fetched.</li>
<li><strong>ID:</strong> Opcode decoded. Registers rs and rt are read
for OR; rs and immediate for ORI.</li>
<li><strong>EX:</strong> ALU performs an OR operation.</li>
<li><strong>MEM:</strong> No action needed.</li>
<li><strong>WB:</strong> Result written to rd (OR) or rt (ORI).</li>
</ul>
<h4 id="instruction-breakdown-2">Instruction Breakdown</h4>
<p>The <code>OR</code> instruction in the MIPS architecture performs a
bitwise OR operation on two register values and stores the result in a
destination register.</p>
<ol type="1">
<li><p><strong>Instruction Fetch (IF) Stage:</strong></p>
<ul>
<li>The Program Counter (PC) holds the address of the <code>OR</code>
instruction to be fetched.</li>
<li>The Instruction Memory module (<code>InstructionMemory.v</code>)
retrieves the instruction from the memory based on the PC value.</li>
<li>The fetched instruction is passed to the next stage.</li>
</ul>
<p>Example code from <code>InstructionMemory.v</code>:
<code>verilog always @(i_Addr) begin   if (i_Addr == -4) begin     i_Instruction = 32'b11111100000000000000000000000000;   end else begin     i_Instruction = Imem[i_Addr&gt;&gt;2];   end   i_Ctr = i_Instruction[31:26];   i_Funcode = i_Instruction[5:0]; end</code></p></li>
<li><p><strong>Instruction Decode (ID) Stage:</strong></p></li>
</ol>
<ul>
<li>The fetched instruction is decoded to identify the opcode and
register operands.</li>
<li>The Control Unit module (<code>ControlUnit.v</code>) sets the
appropriate control signals based on the opcode.</li>
<li>For the <code>OR</code> instruction, the <code>ALUOp</code> control
signal is set to indicate an OR operation.</li>
<li>The register operands (<code>rs</code> and <code>rt</code>) are read
from the Register File.</li>
</ul>
<p>Example code from <code>ControlUnit.v</code>:</p>
<div class="sourceCode" id="cb44"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb44-1"><a href="#cb44-1" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb44-2"><a href="#cb44-2" aria-hidden="true" tabindex="-1"></a>  <span class="kw">case</span> <span class="op">(</span>i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">])</span></span>
<span id="cb44-3"><a href="#cb44-3" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb44-4"><a href="#cb44-4" aria-hidden="true" tabindex="-1"></a>    <span class="bn">6&#39;b000000</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// ARITHMETIC</span></span>
<span id="cb44-5"><a href="#cb44-5" aria-hidden="true" tabindex="-1"></a>      o_RegDst <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb44-6"><a href="#cb44-6" aria-hidden="true" tabindex="-1"></a>      o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-7"><a href="#cb44-7" aria-hidden="true" tabindex="-1"></a>      o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-8"><a href="#cb44-8" aria-hidden="true" tabindex="-1"></a>      o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb44-9"><a href="#cb44-9" aria-hidden="true" tabindex="-1"></a>      o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-10"><a href="#cb44-10" aria-hidden="true" tabindex="-1"></a>      o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-11"><a href="#cb44-11" aria-hidden="true" tabindex="-1"></a>      o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-12"><a href="#cb44-12" aria-hidden="true" tabindex="-1"></a>      o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-13"><a href="#cb44-13" aria-hidden="true" tabindex="-1"></a>      o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b10</span><span class="op">;</span></span>
<span id="cb44-14"><a href="#cb44-14" aria-hidden="true" tabindex="-1"></a>      o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb44-15"><a href="#cb44-15" aria-hidden="true" tabindex="-1"></a>      <span class="co">// ...</span></span>
<span id="cb44-16"><a href="#cb44-16" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb44-17"><a href="#cb44-17" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb44-18"><a href="#cb44-18" aria-hidden="true" tabindex="-1"></a>  <span class="kw">endcase</span></span>
<span id="cb44-19"><a href="#cb44-19" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<ol start="3" type="1">
<li><strong>Execute (EX) Stage:</strong></li>
</ol>
<ul>
<li>The ALU module (<code>ALU.v</code>) performs the bitwise OR
operation on the values of <code>rs</code> and <code>rt</code> based on
the <code>ALUOp</code> control signal.</li>
<li>The result of the OR operation is stored in a temporary
register.</li>
</ul>
<p>Example code from <code>ALU.v</code>:</p>
<div class="sourceCode" id="cb45"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb45-1"><a href="#cb45-1" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> data2<span class="op">,</span> i_ALUcontrol<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb45-2"><a href="#cb45-2" aria-hidden="true" tabindex="-1"></a>  <span class="kw">case</span> <span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb45-3"><a href="#cb45-3" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb45-4"><a href="#cb45-4" aria-hidden="true" tabindex="-1"></a>    <span class="bn">4&#39;b0001</span><span class="op">:</span>  <span class="co">// OR</span></span>
<span id="cb45-5"><a href="#cb45-5" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">|</span> data2<span class="op">;</span> <span class="co">// bitwise OR</span></span>
<span id="cb45-6"><a href="#cb45-6" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb45-7"><a href="#cb45-7" aria-hidden="true" tabindex="-1"></a>  <span class="kw">endcase</span></span>
<span id="cb45-8"><a href="#cb45-8" aria-hidden="true" tabindex="-1"></a>  <span class="co">// ...</span></span>
<span id="cb45-9"><a href="#cb45-9" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<ol start="4" type="1">
<li><strong>Memory Access (MEM) Stage:</strong>
<ul>
<li>For the <code>OR</code> instruction, no memory access is needed, so
this stage is a pass-through.</li>
</ul></li>
<li><strong>Write Back (WB) Stage:</strong></li>
</ol>
<ul>
<li>The ALU result, which is the result of the OR operation, is written
back to the destination register (<code>rd</code>) in the Register
File.</li>
</ul>
<p>Example code for writing back to the Register File:</p>
<div class="sourceCode" id="cb46"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb46-1"><a href="#cb46-1" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb46-2"><a href="#cb46-2" aria-hidden="true" tabindex="-1"></a>  <span class="kw">if</span> <span class="op">(</span>i_RegWrite<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb46-3"><a href="#cb46-3" aria-hidden="true" tabindex="-1"></a>    RegData<span class="op">[</span>i_wReg<span class="op">]</span> <span class="op">=</span> i_wData<span class="op">;</span></span>
<span id="cb46-4"><a href="#cb46-4" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb46-5"><a href="#cb46-5" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<p>Throughout the execution of the <code>OR</code> instruction, the
control signals generated by the Control Unit module
(<code>ControlUnit.v</code>) orchestrate the flow of data and the
operations performed in each stage.</p>
<p>The ALU Control module (<code>ALUControl.v</code>) decodes the
<code>ALUOp</code> signal and the function code of the instruction to
generate the appropriate <code>ALUControl</code> signal for the ALU
module (<code>ALU.v</code>) to perform the OR operation.</p>
<h3 id="ori-or-immediate">ORI (OR Immediate)</h3>
<ul>
<li><strong>IF:</strong> Instruction is fetched.</li>
<li><strong>ID:</strong> Opcode is decoded. Registers rs and rt are read
for OR; rs and immediate for ORI.</li>
<li><strong>EX:</strong> ALU performs an OR operation.</li>
<li><strong>MEM:</strong> No action needed.</li>
<li><strong>WB:</strong> Result is written to rd (OR) or rt (ORI).</li>
</ul>
<h4 id="instruction-breakdown-3">Instruction Breakdown</h4>
<p>The <code>ORI</code> (OR Immediate) instruction in MIPS is an I-type
instruction that performs a bitwise OR operation between a register and
a zero-extended immediate value.</p>
<p>Instruction Format:</p>
<pre><code>ORI rt, rs, immediate</code></pre>
<ul>
<li><code>rt</code>: The destination register where the result will be
stored.</li>
<li><code>rs</code>: The source register containing one of the
operands.</li>
<li><code>immediate</code>: A 16-bit immediate value that will be
zero-extended to 32 bits.</li>
</ul>
<p>Instruction Encoding:</p>
<pre><code>| opcode (6 bits) | rs (5 bits) | rt (5 bits) | immediate (16 bits) |</code></pre>
<p>Processor Stages: 1. Instruction Fetch (IF): - The instruction is
fetched from the instruction memory using the PC (Program Counter). -
The PC is incremented by 4 to point to the next instruction.</p>
<ol start="2" type="1">
<li>Instruction Decode (ID):
<ul>
<li>The instruction is decoded by the control unit.</li>
<li>The register file is accessed to read the values of the source
register <code>rs</code> and the destination register
<code>rt</code>.</li>
<li>The 16-bit immediate value is zero-extended to 32 bits.</li>
</ul></li>
</ol>
<p>Verilog code for the instruction decode stage:</p>
<div class="sourceCode" id="cb49"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb49-1"><a href="#cb49-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Control Unit</span></span>
<span id="cb49-2"><a href="#cb49-2" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb49-3"><a href="#cb49-3" aria-hidden="true" tabindex="-1"></a>  <span class="kw">case</span> <span class="op">(</span>i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">])</span></span>
<span id="cb49-4"><a href="#cb49-4" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb49-5"><a href="#cb49-5" aria-hidden="true" tabindex="-1"></a>    <span class="bn">6&#39;b001101</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// ORI</span></span>
<span id="cb49-6"><a href="#cb49-6" aria-hidden="true" tabindex="-1"></a>      o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb49-7"><a href="#cb49-7" aria-hidden="true" tabindex="-1"></a>      o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb49-8"><a href="#cb49-8" aria-hidden="true" tabindex="-1"></a>      o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb49-9"><a href="#cb49-9" aria-hidden="true" tabindex="-1"></a>      o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb49-10"><a href="#cb49-10" aria-hidden="true" tabindex="-1"></a>      o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb49-11"><a href="#cb49-11" aria-hidden="true" tabindex="-1"></a>      o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb49-12"><a href="#cb49-12" aria-hidden="true" tabindex="-1"></a>      o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb49-13"><a href="#cb49-13" aria-hidden="true" tabindex="-1"></a>      o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b11</span><span class="op">;</span>  <span class="co">// ALU control for OR operation</span></span>
<span id="cb49-14"><a href="#cb49-14" aria-hidden="true" tabindex="-1"></a>      o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb49-15"><a href="#cb49-15" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb49-16"><a href="#cb49-16" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb49-17"><a href="#cb49-17" aria-hidden="true" tabindex="-1"></a>  <span class="kw">endcase</span></span>
<span id="cb49-18"><a href="#cb49-18" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<ol start="3" type="1">
<li>Execute (EX):
<ul>
<li>The ALU performs the bitwise OR operation between the value in the
source register <code>rs</code> and the zero-extended immediate
value.</li>
<li>The ALU control unit generates the appropriate control signal for
the OR operation based on the <code>ALUOp</code> signal from the control
unit.</li>
</ul></li>
</ol>
<p>Verilog code for the ALU control unit:</p>
<div class="sourceCode" id="cb50"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb50-1"><a href="#cb50-1" aria-hidden="true" tabindex="-1"></a><span class="co">// ALU Control</span></span>
<span id="cb50-2"><a href="#cb50-2" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_ALUOp <span class="dt">or</span> i_Funcode<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb50-3"><a href="#cb50-3" aria-hidden="true" tabindex="-1"></a>  <span class="kw">case</span> <span class="op">(</span>i_ALUOp<span class="op">)</span></span>
<span id="cb50-4"><a href="#cb50-4" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb50-5"><a href="#cb50-5" aria-hidden="true" tabindex="-1"></a>    <span class="bn">2&#39;b11</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// ORI</span></span>
<span id="cb50-6"><a href="#cb50-6" aria-hidden="true" tabindex="-1"></a>      o_ALUcontrol <span class="op">=</span> <span class="bn">4&#39;b0001</span><span class="op">;</span>  <span class="co">// ALU control for OR operation</span></span>
<span id="cb50-7"><a href="#cb50-7" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb50-8"><a href="#cb50-8" aria-hidden="true" tabindex="-1"></a>    <span class="co">// ...</span></span>
<span id="cb50-9"><a href="#cb50-9" aria-hidden="true" tabindex="-1"></a>  <span class="kw">endcase</span></span>
<span id="cb50-10"><a href="#cb50-10" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<ol start="4" type="1">
<li>Memory (MEM):
<ul>
<li>No memory access is needed for the <code>ORI</code> instruction, so
this stage is a pass-through.</li>
</ul></li>
<li>Write Back (WB):
<ul>
<li>The ALU result is written back to the destination register
<code>rt</code> in the register file.</li>
</ul></li>
</ol>
<p>Verilog code for the write-back stage:</p>
<div class="sourceCode" id="cb51"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb51-1"><a href="#cb51-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Register File</span></span>
<span id="cb51-2"><a href="#cb51-2" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb51-3"><a href="#cb51-3" aria-hidden="true" tabindex="-1"></a>  <span class="kw">if</span> <span class="op">(</span>i_RegWrite<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb51-4"><a href="#cb51-4" aria-hidden="true" tabindex="-1"></a>    RegData<span class="op">[</span>i_writeReg<span class="op">]</span> <span class="op">&lt;=</span> i_writeData<span class="op">;</span></span>
<span id="cb51-5"><a href="#cb51-5" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb51-6"><a href="#cb51-6" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<p>Here’s an example of how the <code>ORI</code> instruction would be
processed in the single-cycle MIPS processor:</p>
<pre class="assembly"><code>ORI $t0, $s0, 0xFFFF</code></pre>
<p>This instruction performs a bitwise OR operation between the value in
register <code>$s0</code> and the immediate value <code>0xFFFF</code>
(16 bits), and stores the result in register <code>$t0</code>.</p>
<p>In the IF stage, the instruction is fetched from the instruction
memory using the PC.</p>
<p>In the ID stage, the instruction is decoded, and the values of
<code>$s0</code> and <code>$t0</code> are read from the register file.
The immediate value <code>0xFFFF</code> is zero-extended to 32 bits.</p>
<p>In the EX stage, the ALU performs the bitwise OR operation between
the value in <code>$s0</code> and the zero-extended immediate value. The
ALU control unit generates the appropriate control signal for the OR
operation based on the <code>ALUOp</code> signal from the control
unit.</p>
<h3 id="nor">NOR</h3>
<h4 id="intruction-overview">Intruction Overview</h4>
<ul>
<li><strong>IF:</strong> Fetch instruction.</li>
<li><strong>ID:</strong> Decode opcode; read rs and rt.</li>
<li><strong>EX:</strong> ALU performs NOR operation on rs and rt.</li>
<li><strong>MEM:</strong> No action needed.</li>
<li><strong>WB:</strong> Result is written back to rd.</li>
</ul>
<h4 id="instruction-breakdown-4">Instruction Breakdown</h4>
<p>The NOR instruction in the MIPS architecture performs a bitwise NOR
operation on the values of two registers and stores the result in a
destination register.</p>
<ol type="1">
<li><strong>Instruction Fetch (IF) Stage:</strong>
<ul>
<li>The instruction is fetched from the instruction memory using the
program counter (PC).</li>
<li>The instruction memory module (<code>InstructionMemory.v</code>)
retrieves the instruction based on the address provided by the PC.</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb53"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb53-1"><a href="#cb53-1" aria-hidden="true" tabindex="-1"></a> <span class="co">// Inside the InstructionMemory module</span></span>
<span id="cb53-2"><a href="#cb53-2" aria-hidden="true" tabindex="-1"></a> <span class="kw">always</span> <span class="op">@(</span>i_Addr<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb53-3"><a href="#cb53-3" aria-hidden="true" tabindex="-1"></a>   <span class="kw">if</span> <span class="op">(</span>i_Addr <span class="op">==</span> <span class="op">-</span><span class="dv">4</span><span class="op">)</span> <span class="kw">begin</span>         <span class="co">// init</span></span>
<span id="cb53-4"><a href="#cb53-4" aria-hidden="true" tabindex="-1"></a>     i_Instruction <span class="op">=</span> <span class="bn">32&#39;b11111100000000000000000000000000</span><span class="op">;</span></span>
<span id="cb53-5"><a href="#cb53-5" aria-hidden="true" tabindex="-1"></a>   <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb53-6"><a href="#cb53-6" aria-hidden="true" tabindex="-1"></a>     i_Instruction <span class="op">=</span> Imem<span class="op">[</span>i_Addr<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span>
<span id="cb53-7"><a href="#cb53-7" aria-hidden="true" tabindex="-1"></a>   <span class="kw">end</span></span>
<span id="cb53-8"><a href="#cb53-8" aria-hidden="true" tabindex="-1"></a>   i_Ctr <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb53-9"><a href="#cb53-9" aria-hidden="true" tabindex="-1"></a>   i_Funcode <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb53-10"><a href="#cb53-10" aria-hidden="true" tabindex="-1"></a> <span class="kw">end</span></span></code></pre></div>
<ol start="2" type="1">
<li><strong>Instruction Decode (ID) Stage:</strong>
<ul>
<li>The fetched instruction is decoded by the control unit
(<code>ControlUnit.v</code>).</li>
<li>The opcode of the instruction (bits [31:26]) is used to determine
the type of instruction.</li>
<li>For the NOR instruction, the opcode is 6’b100111 (binary
representation of 39).</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb54"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb54-1"><a href="#cb54-1" aria-hidden="true" tabindex="-1"></a> <span class="co">// Inside the ControlUnit module</span></span>
<span id="cb54-2"><a href="#cb54-2" aria-hidden="true" tabindex="-1"></a> <span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb54-3"><a href="#cb54-3" aria-hidden="true" tabindex="-1"></a>   <span class="kw">case</span> <span class="op">(</span>i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">])</span></span>
<span id="cb54-4"><a href="#cb54-4" aria-hidden="true" tabindex="-1"></a>     <span class="co">// ...</span></span>
<span id="cb54-5"><a href="#cb54-5" aria-hidden="true" tabindex="-1"></a>     <span class="bn">6&#39;b100111</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// NOR</span></span>
<span id="cb54-6"><a href="#cb54-6" aria-hidden="true" tabindex="-1"></a>       o_RegDst <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb54-7"><a href="#cb54-7" aria-hidden="true" tabindex="-1"></a>       o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb54-8"><a href="#cb54-8" aria-hidden="true" tabindex="-1"></a>       o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb54-9"><a href="#cb54-9" aria-hidden="true" tabindex="-1"></a>       o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb54-10"><a href="#cb54-10" aria-hidden="true" tabindex="-1"></a>       o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb54-11"><a href="#cb54-11" aria-hidden="true" tabindex="-1"></a>       o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb54-12"><a href="#cb54-12" aria-hidden="true" tabindex="-1"></a>       o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb54-13"><a href="#cb54-13" aria-hidden="true" tabindex="-1"></a>       o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b11</span><span class="op">;</span></span>
<span id="cb54-14"><a href="#cb54-14" aria-hidden="true" tabindex="-1"></a>       o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb54-15"><a href="#cb54-15" aria-hidden="true" tabindex="-1"></a>       <span class="co">// ...</span></span>
<span id="cb54-16"><a href="#cb54-16" aria-hidden="true" tabindex="-1"></a>     <span class="kw">end</span></span>
<span id="cb54-17"><a href="#cb54-17" aria-hidden="true" tabindex="-1"></a>     <span class="co">// ...</span></span>
<span id="cb54-18"><a href="#cb54-18" aria-hidden="true" tabindex="-1"></a>   <span class="kw">endcase</span></span>
<span id="cb54-19"><a href="#cb54-19" aria-hidden="true" tabindex="-1"></a> <span class="kw">end</span></span></code></pre></div>
<ol start="3" type="1">
<li><strong>Execute (EX) Stage:</strong>
<ul>
<li>The ALU performs the bitwise NOR operation on the values of the
source registers (rs and rt).</li>
<li>The ALU control unit generates the appropriate control signal
(4’b1100) for the NOR operation based on the ALUOp bits from the control
unit.</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb55"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb55-1"><a href="#cb55-1" aria-hidden="true" tabindex="-1"></a> <span class="co">// Inside the ALU module</span></span>
<span id="cb55-2"><a href="#cb55-2" aria-hidden="true" tabindex="-1"></a> <span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> data2<span class="op">,</span> i_ALUcontrol<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb55-3"><a href="#cb55-3" aria-hidden="true" tabindex="-1"></a>   <span class="kw">case</span> <span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb55-4"><a href="#cb55-4" aria-hidden="true" tabindex="-1"></a>     <span class="co">// ...</span></span>
<span id="cb55-5"><a href="#cb55-5" aria-hidden="true" tabindex="-1"></a>     <span class="bn">4&#39;b1100</span><span class="op">:</span>  <span class="co">// NOR</span></span>
<span id="cb55-6"><a href="#cb55-6" aria-hidden="true" tabindex="-1"></a>       o_ALUresult <span class="op">=</span> i_data1 <span class="op">|</span> <span class="op">~</span>data2<span class="op">;</span></span>
<span id="cb55-7"><a href="#cb55-7" aria-hidden="true" tabindex="-1"></a>     <span class="co">// ...</span></span>
<span id="cb55-8"><a href="#cb55-8" aria-hidden="true" tabindex="-1"></a>   <span class="kw">endcase</span></span>
<span id="cb55-9"><a href="#cb55-9" aria-hidden="true" tabindex="-1"></a>   <span class="co">// ...</span></span>
<span id="cb55-10"><a href="#cb55-10" aria-hidden="true" tabindex="-1"></a> <span class="kw">end</span></span></code></pre></div>
<ol start="4" type="1">
<li><strong>Memory (MEM) Stage:</strong>
<ul>
<li>For the NOR instruction, no memory access is required, so this stage
is a pass-through.</li>
</ul></li>
<li><strong>Write Back (WB) Stage:</strong>
<ul>
<li>The result of the NOR operation from the ALU is written back to the
destination register (rd) in the register file.</li>
</ul></li>
</ol>
<div class="sourceCode" id="cb56"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb56-1"><a href="#cb56-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Inside the RegisterFile module</span></span>
<span id="cb56-2"><a href="#cb56-2" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb56-3"><a href="#cb56-3" aria-hidden="true" tabindex="-1"></a>  <span class="kw">if</span> <span class="op">(</span>i_wEn <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb56-4"><a href="#cb56-4" aria-hidden="true" tabindex="-1"></a>    RegData<span class="op">[</span>i_wDst<span class="op">]</span> <span class="op">&lt;=</span> i_wData<span class="op">;</span></span>
<span id="cb56-5"><a href="#cb56-5" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb56-6"><a href="#cb56-6" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span></code></pre></div>
<p>Here’s an example of how the NOR instruction would be encoded in MIPS
assembly and its corresponding machine code:</p>
<pre class="assembly"><code># MIPS Assembly
nor $t0, $s1, $s2   # Perform bitwise NOR of $s1 and $s2 and store the result in $t0</code></pre>
<pre><code># Machine Code (in binary)
000000 10001 10010 01000 00000 100111</code></pre>
<p>In the machine code, the bits are organized as follows: - Bits
[31:26]: Opcode (000000 for R-type instructions) - Bits [25:21]: Source
register 1 (rs) - Bits [20:16]: Source register 2 (rt) - Bits [15:11]:
Destination register (rd) BNE (Branch Not Equal)</p>
<ul>
<li><strong>IF:</strong> Fetch instruction.</li>
<li><strong>ID:</strong> Decode instruction; read registers rs and
rt.</li>
<li><strong>EX:</strong> Compare values of rs and rt.</li>
<li><strong>MEM:</strong> No action needed.</li>
<li><strong>WB:</strong> If rs != rt, PC is updated to branch address
(PC + offset); otherwise, move to next sequential instruction.</li>
</ul>
<h4 id="instruction-breakdown-5">Instruction Breakdown</h4>
<p>The BNE (Branch Not Equal) instruction is a conditional branch
instruction in the MIPS architecture.</p>
<p>It compares the values of two registers and transfers control to a
target address if the values are not equal.</p>
<ol type="1">
<li><strong>Instruction Fetch (IF) Stage:</strong> In the IF stage, the
instruction is fetched from the Instruction Memory using the Program
Counter (PC).</li>
</ol>
<div class="sourceCode" id="cb59"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb59-1"><a href="#cb59-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Fetch the instruction from the Instruction Memory</span></span>
<span id="cb59-2"><a href="#cb59-2" aria-hidden="true" tabindex="-1"></a>i_Instruction <span class="op">=</span> Imem<span class="op">[</span>i_Addr<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span></code></pre></div>
<ol start="2" type="1">
<li><strong>Instruction Decode (ID) Stage:</strong> In the ID stage, the
instruction is decoded, and the registers rs and rt are read from the
Register File.</li>
</ol>
<div class="sourceCode" id="cb60"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb60-1"><a href="#cb60-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Decode the instruction</span></span>
<span id="cb60-2"><a href="#cb60-2" aria-hidden="true" tabindex="-1"></a>i_Ctr <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb60-3"><a href="#cb60-3" aria-hidden="true" tabindex="-1"></a>i_Funcode <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb60-4"><a href="#cb60-4" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb60-5"><a href="#cb60-5" aria-hidden="true" tabindex="-1"></a><span class="co">// Read registers rs and rt from the Register File</span></span>
<span id="cb60-6"><a href="#cb60-6" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">25</span><span class="op">:</span><span class="dv">21</span><span class="op">];</span></span>
<span id="cb60-7"><a href="#cb60-7" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> <span class="op">[</span><span class="dv">4</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">20</span><span class="op">:</span><span class="dv">16</span><span class="op">];</span></span>
<span id="cb60-8"><a href="#cb60-8" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs_data <span class="op">=</span> RegData<span class="op">[</span>rs<span class="op">];</span></span>
<span id="cb60-9"><a href="#cb60-9" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt_data <span class="op">=</span> RegData<span class="op">[</span>rt<span class="op">];</span></span></code></pre></div>
<ol start="3" type="1">
<li><strong>Execution (EX) Stage:</strong> In the EX stage, the values
of rs and rt are compared using the ALU.</li>
</ol>
<div class="sourceCode" id="cb61"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb61-1"><a href="#cb61-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Compare the values of rs and rt using the ALU</span></span>
<span id="cb61-2"><a href="#cb61-2" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> ALU_zero<span class="op">;</span></span>
<span id="cb61-3"><a href="#cb61-3" aria-hidden="true" tabindex="-1"></a>ALU alu <span class="op">(</span></span>
<span id="cb61-4"><a href="#cb61-4" aria-hidden="true" tabindex="-1"></a>    .i_data1<span class="op">(</span>rs_data<span class="op">),</span></span>
<span id="cb61-5"><a href="#cb61-5" aria-hidden="true" tabindex="-1"></a>    .i_read2<span class="op">(</span>rt_data<span class="op">),</span></span>
<span id="cb61-6"><a href="#cb61-6" aria-hidden="true" tabindex="-1"></a>    .i_ALUcontrol<span class="op">(</span>ALU_control<span class="op">),</span></span>
<span id="cb61-7"><a href="#cb61-7" aria-hidden="true" tabindex="-1"></a>    .o_Zero<span class="op">(</span>ALU_zero<span class="op">),</span></span>
<span id="cb61-8"><a href="#cb61-8" aria-hidden="true" tabindex="-1"></a>    .o_ALUresult<span class="op">(</span>ALU_result<span class="op">)</span></span>
<span id="cb61-9"><a href="#cb61-9" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span></code></pre></div>
<ol start="4" type="1">
<li><strong>Memory Access (MEM) Stage:</strong></li>
</ol>
<p>The BNE instruction does not require any memory access, so this stage
is a pass-through.</p>
<ol start="5" type="1">
<li><strong>Write Back (WB) Stage:</strong></li>
</ol>
<p>In the WB stage, if the values of rs and rt are not equal (ALU_zero
is 0), the Program Counter (PC) is updated to the branch target address
(PC + offset). Otherwise, the PC moves to the next sequential
instruction.</p>
<div class="sourceCode" id="cb62"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb62-1"><a href="#cb62-1" aria-hidden="true" tabindex="-1"></a><span class="co">// Update the PC based on the branch condition</span></span>
<span id="cb62-2"><a href="#cb62-2" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> branch_target <span class="op">=</span> PC <span class="op">+</span> <span class="op">{{</span><span class="dv">14</span><span class="op">{</span>i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">]}},</span> i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">],</span> <span class="bn">2&#39;b00</span><span class="op">};</span></span>
<span id="cb62-3"><a href="#cb62-3" aria-hidden="true" tabindex="-1"></a><span class="dt">wire</span> branch_taken <span class="op">=</span> i_Branch <span class="op">&amp;</span> <span class="op">~</span>ALU_zero<span class="op">;</span></span>
<span id="cb62-4"><a href="#cb62-4" aria-hidden="true" tabindex="-1"></a><span class="kw">assign</span> PC_next <span class="op">=</span> branch_taken <span class="op">?</span> branch_target <span class="op">:</span> PC <span class="op">+</span> <span class="dv">4</span><span class="op">;</span></span></code></pre></div>
<p>Here’s an example of how the BNE instruction can be represented in
Verilog:</p>
<div class="sourceCode" id="cb63"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb63-1"><a href="#cb63-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> BNE_example <span class="op">(</span></span>
<span id="cb63-2"><a href="#cb63-2" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rs_data<span class="op">,</span></span>
<span id="cb63-3"><a href="#cb63-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> rt_data<span class="op">,</span></span>
<span id="cb63-4"><a href="#cb63-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> offset<span class="op">,</span></span>
<span id="cb63-5"><a href="#cb63-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> PC<span class="op">,</span></span>
<span id="cb63-6"><a href="#cb63-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> PC_next</span>
<span id="cb63-7"><a href="#cb63-7" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb63-8"><a href="#cb63-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> ALU_zero<span class="op">;</span></span>
<span id="cb63-9"><a href="#cb63-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">wire</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> branch_target <span class="op">=</span> PC <span class="op">+</span> <span class="op">{{</span><span class="dv">14</span><span class="op">{</span>offset<span class="op">[</span><span class="dv">15</span><span class="op">]}},</span> offset<span class="op">,</span> <span class="bn">2&#39;b00</span><span class="op">};</span></span>
<span id="cb63-10"><a href="#cb63-10" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Compare rs and rt using ALU</span></span>
<span id="cb63-11"><a href="#cb63-11" aria-hidden="true" tabindex="-1"></a>    <span class="kw">assign</span> ALU_zero <span class="op">=</span> <span class="op">(</span>rs_data <span class="op">==</span> rt_data<span class="op">)</span> <span class="op">?</span> <span class="dv">1</span> <span class="op">:</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb63-12"><a href="#cb63-12" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Update PC based on branch condition</span></span>
<span id="cb63-13"><a href="#cb63-13" aria-hidden="true" tabindex="-1"></a>    <span class="kw">always</span> <span class="op">@(</span>*<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb63-14"><a href="#cb63-14" aria-hidden="true" tabindex="-1"></a>        <span class="kw">if</span> <span class="op">(~</span>ALU_zero<span class="op">)</span></span>
<span id="cb63-15"><a href="#cb63-15" aria-hidden="true" tabindex="-1"></a>            PC_next <span class="op">=</span> branch_target<span class="op">;</span></span>
<span id="cb63-16"><a href="#cb63-16" aria-hidden="true" tabindex="-1"></a>        <span class="kw">else</span></span>
<span id="cb63-17"><a href="#cb63-17" aria-hidden="true" tabindex="-1"></a>            PC_next <span class="op">=</span> PC <span class="op">+</span> <span class="dv">4</span><span class="op">;</span></span>
<span id="cb63-18"><a href="#cb63-18" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb63-19"><a href="#cb63-19" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>In this example, the <code>BNE_example</code> module takes the values
of rs and rt (<code>rs_data</code> and <code>rt_data</code>), the branch
offset (<code>offset</code>), and the current Program Counter
(<code>PC</code>) as inputs.</p>
<p>It compares rs and rt using the ALU and updates the
<code>PC_next</code> output based on the branch condition.</p>
<p>If rs and rt are not equal (<code>ALU_zero</code> is 0),
<code>PC_next</code> is set to the branch target address
(<code>branch_target</code>).</p>
<p>Otherwise, <code>PC_next</code> is set to the next sequential
instruction address (<code>PC + 4</code>).</p>
<h1 id="comparing-verilog-vs-vhdl">Comparing Verilog vs VHDL</h1>
<p>The following section will compare the experience of developing a
single-cycle processor in Verilog vs VHDL.</p>
<h2 id="interesting-notes-about-verilog">Interesting notes about
verilog</h2>
<p>The loose typing of Verilog can lead to some useful modules that can
be defined in small amounts of code.</p>
<p>Additionally, in verilog, you do not have to declare your components
in the component that uses them which also decreases the amount of code
that is needed to be written.</p>
<p>For example, the following module is the mux module that is used in
the processor to select between two inputs based on a control
signal.</p>
<div class="sourceCode" id="cb64"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb64-1"><a href="#cb64-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> mux #<span class="op">(</span><span class="dt">parameter</span> size <span class="op">=</span> <span class="dv">1</span><span class="op">)</span> <span class="op">(</span></span>
<span id="cb64-2"><a href="#cb64-2" aria-hidden="true" tabindex="-1"></a>  <span class="dt">input</span> select<span class="op">,</span></span>
<span id="cb64-3"><a href="#cb64-3" aria-hidden="true" tabindex="-1"></a>  <span class="dt">input</span> <span class="op">[</span>size <span class="op">-</span> <span class="dv">1</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> in_0<span class="op">,</span></span>
<span id="cb64-4"><a href="#cb64-4" aria-hidden="true" tabindex="-1"></a>  <span class="dt">input</span> <span class="op">[</span>size <span class="op">-</span> <span class="dv">1</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> in_1<span class="op">,</span></span>
<span id="cb64-5"><a href="#cb64-5" aria-hidden="true" tabindex="-1"></a>  <span class="dt">output</span> <span class="op">[</span>size <span class="op">-</span> <span class="dv">1</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> out</span>
<span id="cb64-6"><a href="#cb64-6" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb64-7"><a href="#cb64-7" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> out <span class="op">=</span> <span class="op">(</span>select<span class="op">)</span> <span class="op">?</span> in_1 <span class="op">:</span> in_0<span class="op">;</span></span>
<span id="cb64-8"><a href="#cb64-8" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Another example of this is the sign extender module that is used to
extend the sign of a 16-bit number to a 32-bit number.</p>
<div class="sourceCode" id="cb65"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb65-1"><a href="#cb65-1" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> signextender <span class="op">(</span></span>
<span id="cb65-2"><a href="#cb65-2" aria-hidden="true" tabindex="-1"></a>  <span class="dt">input</span> <span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> in<span class="op">,</span></span>
<span id="cb65-3"><a href="#cb65-3" aria-hidden="true" tabindex="-1"></a>  <span class="dt">output</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> out</span>
<span id="cb65-4"><a href="#cb65-4" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb65-5"><a href="#cb65-5" aria-hidden="true" tabindex="-1"></a>  <span class="kw">assign</span> out <span class="op">=</span> <span class="op">{{</span><span class="dv">16</span><span class="op">{</span>in<span class="op">[</span><span class="dv">15</span><span class="op">]}},</span> <span class="op">{</span>in<span class="op">}};</span></span>
<span id="cb65-6"><a href="#cb65-6" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Additionally, the fact that in Verilog you do not need to
preemptively define components before using them allows for a more
flexible design and faster development.</p>
<p>Another nice feature of Verilog (specifically VerilogHDL) is the
ability to print out the values of signals without having to worry about
their types in the waveform viewer inside modelsim/questasim.</p>
<p>This is a feature that is not present in VHDL and is very useful for
debugging and understanding the behavior of the processor.</p>
<h2 id="interesting-notes-about-vhdl">Interesting notes about VHDL</h2>
<p>While you can do this “print debugging” in VHDL, it is not as easy as
it is in Verilog because in VHDL you must deal with the typings of the
signals and the fact that you must declare the signals before you can
use them.</p>
<p>As an example, the following is the code in VHDL that prints out the
values of the signals in the waveform viewer for an n-bit register which
was used in a project for CPRE381.</p>
<p>The following test-bench shows the code that can be executed inside
modelsim/questasim to print out the values of the signals in the
waveform viewer.</p>
<p>It displays the additional hurdles that must be overcome in VHDL to
print out the values of the signals in the waveform viewer because of
the strict typing of the language.</p>
<div class="sourceCode" id="cb66"><pre
class="sourceCode vhdl"><code class="sourceCode vhdl"><span id="cb66-1"><a href="#cb66-1" aria-hidden="true" tabindex="-1"></a><span class="kw">LIBRARY</span> IEEE;</span>
<span id="cb66-2"><a href="#cb66-2" aria-hidden="true" tabindex="-1"></a><span class="kw">USE</span> IEEE<span class="ot">.</span>std_logic_1164<span class="ot">.</span><span class="kw">ALL</span>;</span>
<span id="cb66-3"><a href="#cb66-3" aria-hidden="true" tabindex="-1"></a><span class="cf">ENTITY</span> <span class="dv">tb_nbitregister</span> <span class="kw">IS</span></span>
<span id="cb66-4"><a href="#cb66-4" aria-hidden="true" tabindex="-1"></a>  <span class="cf">GENERIC</span> (</span>
<span id="cb66-5"><a href="#cb66-5" aria-hidden="true" tabindex="-1"></a>    gCLK_HPER <span class="ot">:</span> <span class="dt">TIME</span> <span class="ot">:=</span> <span class="dv">50</span> <span class="dt">ns</span>;</span>
<span id="cb66-6"><a href="#cb66-6" aria-hidden="true" tabindex="-1"></a>    N <span class="ot">:</span> <span class="dt">INTEGER</span> <span class="ot">:=</span> <span class="dv">32</span>);</span>
<span id="cb66-7"><a href="#cb66-7" aria-hidden="true" tabindex="-1"></a><span class="cf">END tb_nbitregister;</span></span>
<span id="cb66-8"><a href="#cb66-8" aria-hidden="true" tabindex="-1"></a><span class="cf">ARCHITECTURE</span> <span class="dv">behavior</span> <span class="kw">OF</span> <span class="fu">tb_nbitregister</span> <span class="kw">IS</span></span>
<span id="cb66-9"><a href="#cb66-9" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb66-10"><a href="#cb66-10" aria-hidden="true" tabindex="-1"></a>  <span class="co">-- Calculate the clock period as twice the half-period</span></span>
<span id="cb66-11"><a href="#cb66-11" aria-hidden="true" tabindex="-1"></a>  <span class="ot">CONSTANT</span> cCLK_PER <span class="ot">:</span> <span class="dt">TIME</span> <span class="ot">:=</span> gCLK_HPER <span class="ot">*</span> <span class="dv">2</span>;</span>
<span id="cb66-12"><a href="#cb66-12" aria-hidden="true" tabindex="-1"></a>  <span class="cf">COMPONENT</span> <span class="dv">nbitregister</span></span>
<span id="cb66-13"><a href="#cb66-13" aria-hidden="true" tabindex="-1"></a>    <span class="cf">PORT</span> (</span>
<span id="cb66-14"><a href="#cb66-14" aria-hidden="true" tabindex="-1"></a>      i_CLK <span class="ot">:</span> <span class="kw">IN</span> <span class="dt">STD_LOGIC</span>; <span class="co">-- Clock input</span></span>
<span id="cb66-15"><a href="#cb66-15" aria-hidden="true" tabindex="-1"></a>      i_RST <span class="ot">:</span> <span class="kw">IN</span> <span class="dt">STD_LOGIC</span>; <span class="co">-- Reset input</span></span>
<span id="cb66-16"><a href="#cb66-16" aria-hidden="true" tabindex="-1"></a>      i_WE <span class="ot">:</span> <span class="kw">IN</span> <span class="dt">STD_LOGIC</span>; <span class="co">-- Write enable input</span></span>
<span id="cb66-17"><a href="#cb66-17" aria-hidden="true" tabindex="-1"></a>      i_D <span class="ot">:</span> <span class="kw">IN</span> <span class="dt">STD_LOGIC_VECTOR</span>(N <span class="ot">-</span> <span class="dv">1</span> <span class="ot">DOWNTO</span> <span class="dv">0</span>); <span class="co">-- Data value input</span></span>
<span id="cb66-18"><a href="#cb66-18" aria-hidden="true" tabindex="-1"></a>      o_Q <span class="ot">:</span> <span class="kw">OUT</span> <span class="dt">STD_LOGIC_VECTOR</span>(N <span class="ot">-</span> <span class="dv">1</span> <span class="ot">DOWNTO</span> <span class="dv">0</span>) <span class="co">-- Data value output</span></span>
<span id="cb66-19"><a href="#cb66-19" aria-hidden="true" tabindex="-1"></a>    );</span>
<span id="cb66-20"><a href="#cb66-20" aria-hidden="true" tabindex="-1"></a>  <span class="cf">END COMPONENT;</span></span>
<span id="cb66-21"><a href="#cb66-21" aria-hidden="true" tabindex="-1"></a>  <span class="co">-- Temporary Signals to connect to the nbitregister component.</span></span>
<span id="cb66-22"><a href="#cb66-22" aria-hidden="true" tabindex="-1"></a>  <span class="ot">SIGNAL</span> s_CLK<span class="ot">,</span> s_RST<span class="ot">,</span> s_WE <span class="ot">:</span> <span class="dt">STD_LOGIC</span>;</span>
<span id="cb66-23"><a href="#cb66-23" aria-hidden="true" tabindex="-1"></a>  <span class="ot">SIGNAL</span> s_D<span class="ot">,</span> s_Q <span class="ot">:</span> <span class="dt">STD_LOGIC_VECTOR</span>(N <span class="ot">-</span> <span class="dv">1</span> <span class="ot">DOWNTO</span> <span class="dv">0</span>);</span>
<span id="cb66-24"><a href="#cb66-24" aria-hidden="true" tabindex="-1"></a><span class="cf">BEGIN</span></span>
<span id="cb66-25"><a href="#cb66-25" aria-hidden="true" tabindex="-1"></a>  <span class="dv">DUT</span> <span class="ot">:</span> <span class="fu">nbitregister</span></span>
<span id="cb66-26"><a href="#cb66-26" aria-hidden="true" tabindex="-1"></a>  <span class="kw">PORT MAP</span>(</span>
<span id="cb66-27"><a href="#cb66-27" aria-hidden="true" tabindex="-1"></a>    i_CLK <span class="ot">=&gt;</span> s_CLK<span class="ot">,</span></span>
<span id="cb66-28"><a href="#cb66-28" aria-hidden="true" tabindex="-1"></a>    i_RST <span class="ot">=&gt;</span> s_RST<span class="ot">,</span></span>
<span id="cb66-29"><a href="#cb66-29" aria-hidden="true" tabindex="-1"></a>    i_WE <span class="ot">=&gt;</span> s_WE<span class="ot">,</span></span>
<span id="cb66-30"><a href="#cb66-30" aria-hidden="true" tabindex="-1"></a>    i_D <span class="ot">=&gt;</span> s_D<span class="ot">,</span></span>
<span id="cb66-31"><a href="#cb66-31" aria-hidden="true" tabindex="-1"></a>    o_Q <span class="ot">=&gt;</span> s_Q</span>
<span id="cb66-32"><a href="#cb66-32" aria-hidden="true" tabindex="-1"></a>  );</span>
<span id="cb66-33"><a href="#cb66-33" aria-hidden="true" tabindex="-1"></a>  <span class="dv">P_CLK</span> <span class="ot">:</span> <span class="sc">PROCESS</span></span>
<span id="cb66-34"><a href="#cb66-34" aria-hidden="true" tabindex="-1"></a>  <span class="sc">BEGIN</span></span>
<span id="cb66-35"><a href="#cb66-35" aria-hidden="true" tabindex="-1"></a>    s_CLK <span class="ot">&lt;=</span> <span class="bn">&#39;0&#39;</span>;</span>
<span id="cb66-36"><a href="#cb66-36" aria-hidden="true" tabindex="-1"></a>    <span class="kw">WAIT</span> <span class="kw">FOR</span> gCLK_HPER;</span>
<span id="cb66-37"><a href="#cb66-37" aria-hidden="true" tabindex="-1"></a>    s_CLK <span class="ot">&lt;=</span> <span class="bn">&#39;1&#39;</span>;</span>
<span id="cb66-38"><a href="#cb66-38" aria-hidden="true" tabindex="-1"></a>    <span class="kw">WAIT</span> <span class="kw">FOR</span> gCLK_HPER;</span>
<span id="cb66-39"><a href="#cb66-39" aria-hidden="true" tabindex="-1"></a>  <span class="sc">END PROCESS;</span></span>
<span id="cb66-40"><a href="#cb66-40" aria-hidden="true" tabindex="-1"></a>  <span class="dv">P_TB</span> <span class="ot">:</span> <span class="sc">PROCESS</span></span>
<span id="cb66-41"><a href="#cb66-41" aria-hidden="true" tabindex="-1"></a>  <span class="sc">BEGIN</span></span>
<span id="cb66-42"><a href="#cb66-42" aria-hidden="true" tabindex="-1"></a>    s_RST <span class="ot">&lt;=</span> <span class="bn">&#39;1&#39;</span>;</span>
<span id="cb66-43"><a href="#cb66-43" aria-hidden="true" tabindex="-1"></a>    s_WE <span class="ot">&lt;=</span> <span class="bn">&#39;0&#39;</span>;</span>
<span id="cb66-44"><a href="#cb66-44" aria-hidden="true" tabindex="-1"></a>    s_D <span class="ot">&lt;=</span> <span class="st">&quot;00000000000000000000000000000000&quot;</span>;</span>
<span id="cb66-45"><a href="#cb66-45" aria-hidden="true" tabindex="-1"></a>    <span class="kw">WAIT</span> <span class="kw">FOR</span> cCLK_PER;</span>
<span id="cb66-46"><a href="#cb66-46" aria-hidden="true" tabindex="-1"></a>    <span class="co">-- </span><span class="al">TEST</span><span class="co"> CASE 1 - STORE &#39;1&#39;</span></span>
<span id="cb66-47"><a href="#cb66-47" aria-hidden="true" tabindex="-1"></a>    <span class="co">-- DESCRIPTION: The register should store the new data value</span></span>
<span id="cb66-48"><a href="#cb66-48" aria-hidden="true" tabindex="-1"></a>    <span class="co">-- EXPECTED RESULT: The new data value should be stored in the register</span></span>
<span id="cb66-49"><a href="#cb66-49" aria-hidden="true" tabindex="-1"></a>    s_RST <span class="ot">&lt;=</span> <span class="bn">&#39;0&#39;</span>;</span>
<span id="cb66-50"><a href="#cb66-50" aria-hidden="true" tabindex="-1"></a>    s_WE <span class="ot">&lt;=</span> <span class="bn">&#39;1&#39;</span>;</span>
<span id="cb66-51"><a href="#cb66-51" aria-hidden="true" tabindex="-1"></a>    s_D <span class="ot">&lt;=</span> <span class="st">&quot;11111111111111111111111111111111&quot;</span>;</span>
<span id="cb66-52"><a href="#cb66-52" aria-hidden="true" tabindex="-1"></a>    <span class="kw">WAIT</span> <span class="kw">FOR</span> cCLK_PER;</span>
<span id="cb66-53"><a href="#cb66-53" aria-hidden="true" tabindex="-1"></a>    <span class="cf">IF</span> (s_Q <span class="ot">/=</span> <span class="st">&quot;11111111111111111111111111111111&quot;</span>) <span class="cf">THEN</span></span>
<span id="cb66-54"><a href="#cb66-54" aria-hidden="true" tabindex="-1"></a>      <span class="kw">REPORT</span> <span class="st">&quot;Test 1 failed&quot;</span>;</span>
<span id="cb66-55"><a href="#cb66-55" aria-hidden="true" tabindex="-1"></a>      <span class="kw">REPORT</span> <span class="st">&quot;Expected: 11111111111111111111111111111111&quot;</span>;</span>
<span id="cb66-56"><a href="#cb66-56" aria-hidden="true" tabindex="-1"></a>      <span class="kw">REPORT</span> <span class="st">&quot;Actual:  &quot;</span> <span class="ot">&amp;</span> STD_LOGIC_VECTOR<span class="bn">&#39;image(s_Q)</span>;</span>
<span id="cb66-57"><a href="#cb66-57" aria-hidden="true" tabindex="-1"></a>    <span class="cf">ELSE</span></span>
<span id="cb66-58"><a href="#cb66-58" aria-hidden="true" tabindex="-1"></a>      <span class="kw">REPORT</span> <span class="st">&quot;TEST 1 PASSED - STORE &#39;1&#39;&quot;</span>;</span>
<span id="cb66-59"><a href="#cb66-59" aria-hidden="true" tabindex="-1"></a>    <span class="cf">END IF;</span></span>
<span id="cb66-60"><a href="#cb66-60" aria-hidden="true" tabindex="-1"></a>      <span class="ot">//</span> <span class="ot">...</span></span>
<span id="cb66-61"><a href="#cb66-61" aria-hidden="true" tabindex="-1"></a><span class="sc">END behavior;</span></span></code></pre></div>
<h2 id="conclusion">Conclusion</h2>
<p>I think that VHDL actually provides more flexibility within the
development of the processor.</p>
<p>While the language is more verbose because you must reinstantiate a
component within another component to use it, it is more type-safe, and
allows for more control over the design of the processor.</p>
<p>Verilog is more concise and easier to read, but I think that VHDL is
more powerful and allows for more control over the design of the
processor because of it’s type-safety.</p>
<p>To support this, I present the result of the actual lines of code
that were written for the processor in VHDL and Verilog.</p>
<h4 id="vhdl-single-cycle-mips-processor-code-statistics">VHDL Single
Cycle MIPS Processor Code Statistics:</h4>
<pre><code>Language                         files          blank        comment           code
-----------------------------------------------------------------------------------
VHDL                                65            694           1085           4677</code></pre>
<h4 id="verilog-single-cycle-mips-processor-code-statistics">Verilog
Single Cycle MIPS Processor Code Statistics:</h4>
<pre><code>Language                         files          blank        comment           code
-----------------------------------------------------------------------------------
Verilog-SystemVerilog               10              0              9            555</code></pre>
<p>As you can see, the VHDL processor has almost <strong>10</strong>
times the amount of code as the Verilog processor.</p>
<p>Furthermore, I think that the fact that the name of a file in verilog
must match the module name is a limitation that VHDL does not have
(atleast in our Quartus simulator).</p>
<p>To summarize, I think that VHDL is more suited for larger projects
and more complex designs where the type-safety and , while Verilog is
more suited for smaller projects and simpler designs.</p>
<h1 id="breaking-down-decoding-a-signal-to-7-segment-displays">Breaking
down decoding a signal to 7-segment displays</h1>
<p>As the signal representing the instruction is 5 bits long inside the
<code>controller.v</code> file, we need to decode this signal to display
the current instruction on the 7-segment displays.</p>
<p>This means that we need to decode a 5-bit signal to a 35-bit signal
that will be used to display the current instruction on the 7-segment
displays.</p>
<p>If 5-bits are used to represent the instruction, and 7-bits are
needed to represent a character on a 7-segment display, then 35-bits are
needed to represent the current instruction on 5 7-segment displays as 7
* 5 = 35.</p>
<p>Furthermore, the longest word that can be displayed on the 7-segment
displays is 5 characters long, so 5 * 7 = 35 bits are needed to
represent the current instruction on the 7-segment displays.</p>
<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: center;">Func_in</th>
<th style="text-align: center;">O_out</th>
<th style="text-align: center;">Operation</th>
<th style="text-align: center;">Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: center;">1000</td>
<td style="text-align: center;">ox</td>
<td style="text-align: center;"><span
class="math inline">(<em>A</em>+<em>B</em>)</span></td>
<td style="text-align: center;">ADD</td>
</tr>
<tr class="even">
<td style="text-align: center;">1000</td>
<td style="text-align: center;"><span class="math inline">1X</span></td>
<td style="text-align: center;"><span
class="math inline">(A−B)</span></td>
<td style="text-align: center;">SuB</td>
</tr>
<tr class="odd">
<td style="text-align: center;">1001</td>
<td style="text-align: center;">00</td>
<td style="text-align: center;"><span
class="math inline">(<em>A</em>&amp;<em>B</em>)</span></td>
<td style="text-align: center;">AND</td>
</tr>
<tr class="even">
<td style="text-align: center;">1001</td>
<td style="text-align: center;">01</td>
<td style="text-align: center;"><span
class="math inline">(<em>A</em>∣<em>B</em>)</span></td>
<td style="text-align: center;">OR</td>
</tr>
<tr class="odd">
<td style="text-align: center;">1001</td>
<td style="text-align: center;"><span
class="math inline"><em>π</em></span></td>
<td style="text-align: center;"><span
class="math inline"> ∼ (A∣B)</span></td>
<td style="text-align: center;">NOR</td>
</tr>
<tr class="even">
<td style="text-align: center;">101</td>
<td style="text-align: center;"><span
class="math inline">xx0</span></td>
<td style="text-align: center;">signed <span
class="math inline">(<em>A</em>) &lt; signed (<em>B</em>)</span></td>
<td style="text-align: center;">Set-Less-Than signed</td>
</tr>
<tr class="odd">
<td style="text-align: center;">101</td>
<td style="text-align: center;"><span
class="math inline"><em>x</em><em>x</em>1</span></td>
<td style="text-align: center;"><span
class="math inline"><em>A</em> &lt; <em>B</em></span></td>
<td style="text-align: center;">Set-Less-Than unsigned</td>
</tr>
<tr class="even">
<td style="text-align: center;">111</td>
<td style="text-align: center;">000</td>
<td style="text-align: center;">A</td>
<td style="text-align: center;">BLTZ (Branch if Less Than Zero)</td>
</tr>
<tr class="odd">
<td style="text-align: center;">111</td>
<td style="text-align: center;">001</td>
<td style="text-align: center;">A</td>
<td style="text-align: center;">BGEZ (Branch if Greater or Equal to
Zero)</td>
</tr>
<tr class="even">
<td style="text-align: center;">111</td>
<td style="text-align: center;">010</td>
<td style="text-align: center;">A</td>
<td style="text-align: center;">J/AL (Jump and Link)</td>
</tr>
<tr class="odd">
<td style="text-align: center;">111</td>
<td style="text-align: center;">011</td>
<td style="text-align: center;">A</td>
<td style="text-align: center;">JR/AL (Jump Register and Link)</td>
</tr>
<tr class="even">
<td style="text-align: center;">111</td>
<td style="text-align: center;">100</td>
<td style="text-align: center;">A</td>
<td style="text-align: center;"><span class="math inline">BEQ(</span>
(Branch if Equal)</td>
</tr>
<tr class="odd">
<td style="text-align: center;">111</td>
<td style="text-align: center;">101</td>
<td style="text-align: center;">A</td>
<td style="text-align: center;">BNE (Branch if Not Equal)</td>
</tr>
</tbody>
</table>
<p>The following is the wave-diagram from modelsim/questasim for my
test-bench of my processor without the added seven segment displays.</p>
<figure>
<img src="assets/WaveDiagramWithoutSevenSegment.png"
alt="WaveDiagramWithoutSevenSegment.png" />
<figcaption
aria-hidden="true">WaveDiagramWithoutSevenSegment.png</figcaption>
</figure>
<p>Below is the captured wave-diagram from modelsim/questasim with the
seven segment ports included:</p>
<figure>
<img src="assets/SevenSegmentWaveDiagram.png"
alt="SevenSegmentWaveDiagram.png" />
<figcaption aria-hidden="true">SevenSegmentWaveDiagram.png</figcaption>
</figure>
<p>The wave-forms show the output of the following test-bench,
<code>mips_tb.v</code>, which is used to test the single-cycle MIPS
processor, <code>mips.v</code>.</p>
<div class="sourceCode" id="cb69"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb69-1"><a href="#cb69-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb69-2"><a href="#cb69-2" aria-hidden="true" tabindex="-1"></a><span class="ot">`define CYCLE_TIME 20</span></span>
<span id="cb69-3"><a href="#cb69-3" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> mips_tb<span class="op">;</span></span>
<span id="cb69-4"><a href="#cb69-4" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> clk<span class="op">;</span></span>
<span id="cb69-5"><a href="#cb69-5" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> rst<span class="op">;</span></span>
<span id="cb69-6"><a href="#cb69-6" aria-hidden="true" tabindex="-1"></a>  <span class="co">// segments for the 7-segment displays</span></span>
<span id="cb69-7"><a href="#cb69-7" aria-hidden="true" tabindex="-1"></a>  <span class="dt">wire</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> seg_first<span class="op">,</span> seg_second<span class="op">,</span> seg_third<span class="op">,</span> seg_fourth<span class="op">,</span> seg_fifth<span class="op">;</span></span>
<span id="cb69-8"><a href="#cb69-8" aria-hidden="true" tabindex="-1"></a>  <span class="dt">integer</span> i<span class="op">;</span></span>
<span id="cb69-9"><a href="#cb69-9" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> #<span class="op">(</span><span class="ot">`CYCLE_TIME</span> <span class="op">/</span> <span class="dv">2</span><span class="op">)</span> clk <span class="op">=</span> <span class="op">~</span>clk<span class="op">;</span></span>
<span id="cb69-10"><a href="#cb69-10" aria-hidden="true" tabindex="-1"></a>  mips uut <span class="op">(</span></span>
<span id="cb69-11"><a href="#cb69-11" aria-hidden="true" tabindex="-1"></a>      .i_Clk<span class="op">(</span>clk<span class="op">),</span></span>
<span id="cb69-12"><a href="#cb69-12" aria-hidden="true" tabindex="-1"></a>      .i_Rst<span class="op">(</span>rst<span class="op">),</span></span>
<span id="cb69-13"><a href="#cb69-13" aria-hidden="true" tabindex="-1"></a>      .o_Seg_first<span class="op">(</span>seg_first<span class="op">),</span></span>
<span id="cb69-14"><a href="#cb69-14" aria-hidden="true" tabindex="-1"></a>      .o_Seg_second<span class="op">(</span>seg_second<span class="op">),</span></span>
<span id="cb69-15"><a href="#cb69-15" aria-hidden="true" tabindex="-1"></a>      .o_Seg_third<span class="op">(</span>seg_third<span class="op">),</span></span>
<span id="cb69-16"><a href="#cb69-16" aria-hidden="true" tabindex="-1"></a>      .o_Seg_fourth<span class="op">(</span>seg_fourth<span class="op">),</span></span>
<span id="cb69-17"><a href="#cb69-17" aria-hidden="true" tabindex="-1"></a>      .o_Seg_fifth<span class="op">(</span>seg_fifth<span class="op">)</span></span>
<span id="cb69-18"><a href="#cb69-18" aria-hidden="true" tabindex="-1"></a>  <span class="op">);</span></span>
<span id="cb69-19"><a href="#cb69-19" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb69-20"><a href="#cb69-20" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize data memory</span></span>
<span id="cb69-21"><a href="#cb69-21" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb69-22"><a href="#cb69-22" aria-hidden="true" tabindex="-1"></a>      uut.inst_DataMemory.Dmem<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb69-23"><a href="#cb69-23" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb69-24"><a href="#cb69-24" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize Register File</span></span>
<span id="cb69-25"><a href="#cb69-25" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb69-26"><a href="#cb69-26" aria-hidden="true" tabindex="-1"></a>      uut.inst_RegisterFile.RegData<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb69-27"><a href="#cb69-27" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb69-28"><a href="#cb69-28" aria-hidden="true" tabindex="-1"></a>    clk <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb69-29"><a href="#cb69-29" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb69-30"><a href="#cb69-30" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb69-31"><a href="#cb69-31" aria-hidden="true" tabindex="-1"></a>    <span class="bn">#1800</span> <span class="dt">$finish</span><span class="op">;</span></span>
<span id="cb69-32"><a href="#cb69-32" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb69-33"><a href="#cb69-33" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>The given Verilog code represents my test-bench module that was used
for testing the single-cycle MIPS processor, <code>mips.v</code>.</p>
<ol type="1">
<li><p>The test-bench module is named <code>mips_tb</code> (titled
<code>mips_tb.v</code>), and it operates on a timescale of
1ns/1ps.</p></li>
<li><p>The module declares two reg variables:</p>
<ul>
<li><code>clk</code>: Represents the clock signal for the
processor.</li>
<li><code>rst</code>: Represents the reset signal for the
processor.</li>
</ul></li>
<li><p>It also declares five wire variables (<code>seg_first</code>,
<code>seg_second</code>, <code>seg_third</code>,
<code>seg_fourth</code>, <code>seg_fifth</code>) to represent the
segments for the 7-segment displays. These wires are used to display the
current instruction being executed by the processor.</p></li>
<li><p>The <code>integer</code> variable <code>i</code> is declared as a
loop variable for initializing memory.</p></li>
<li><p>The <code>always</code> block generates the clock signal by
toggling the <code>clk</code> variable every half of the clock cycle
time (<code>CYCLE_TIME/2</code>).</p></li>
<li><p>The <code>mips</code> module (the actual MIPS processor) is
instantiated as <code>uut</code> (unit under test) with the following
connections:</p></li>
</ol>
<ul>
<li><code>i_Clk</code> is connected to the <code>clk</code> signal.</li>
<li><code>i_Rst</code> is connected to the <code>rst</code> signal.</li>
<li>The 7-segment display outputs (<code>o_Seg_first</code>,
<code>o_Seg_second</code>, <code>o_Seg_third</code>,
<code>o_Seg_fourth</code>, <code>o_Seg_fifth</code>) are connected to
the corresponding wires in the test-bench.</li>
</ul>
<ol start="7" type="1">
<li>The first <code>initial</code> block is used to initialize the data
memory and the register file of the MIPS processor:
<ul>
<li>It uses a <code>for</code> loop to iterate over the first 32
locations of the data memory (<code>Dmem</code>) and initializes each
location to zero.</li>
<li>Similarly, it initializes the first 32 registers in the register
file (<code>RegData</code>) to zero.</li>
<li>Finally, it sets the <code>clk</code> variable to 0.</li>
</ul></li>
<li>The second <code>initial</code> block is used to specify the
duration of the simulation. It uses the <code>$finish</code> system task
to terminate the simulation after 1800 time units.</li>
</ol>
<p>The purpose of this test-bench is to provide a simulation environment
for the MIPS processor. It initializes the necessary components (data
memory and register file), generates the clock signal, and instantiates
the MIPS processor module. The testbench also specifies the duration of
the simulation.</p>
<p>The test-bench interacts with other components of the processor
through the instantiated <code>mips</code> module (<code>uut</code>). It
provides the clock and reset signals to the processor and observes the
output signals for the 7-segment displays.</p>
<p>Overall, this test-bench serves as a framework to verify the
functionality of the single-cycle MIPS processor by providing the
necessary inputs, initializing the memory, and specifying the simulation
duration.</p>
<h2 id="schematics">Schematics</h2>
<p>The following section will include the schematics for the various
components of the MIPS processor as written in verilog.</p>
<h3 id="control-unit-schematic">Control Unit Schematic</h3>
<p>The following is the schematic for the control unit of the MIPS
processor.</p>
<p>[[Pasted image 20240503123506.png]]</p>
<h3 id="register-file">Register File</h3>
<p>The following is the schematic for the register file of the MIPS
processor.</p>
<p>[[RegisterFile.png]]</p>
<h3 id="data-memory">Data Memory</h3>
<p>The following is the schematic for the data memory of the MIPS
processor.</p>
<figure>
<img src="assets/SchematicDataMemory.png" alt="SchematicDataMemory.png" />
<figcaption aria-hidden="true">SchematicDataMemory.png</figcaption>
</figure>
<h3 id="alu-control">ALU Control</h3>
<p>The following is the schematic for the ALU control of the MIPS
processor.</p>
<figure>
<img src="assets/SchematicALUControl.png" alt="SchematicALUControl.png" />
<figcaption aria-hidden="true">SchematicALUControl.png</figcaption>
</figure>
<h3 id="program-counter-control">Program Counter Control</h3>
<p>The following is the schematic for the program counter control of the
MIPS processor.</p>
<figure>
<img src="assets/SchematicNextProgramCounter.png"
alt="SchematicNextProgramCounter.png" />
<figcaption
aria-hidden="true">SchematicNextProgramCounter.png</figcaption>
</figure>
<h3 id="alu">ALU</h3>
<p>The following is the schematic for the ALU of the MIPS processor.</p>
<figure>
<img src="assets/SchematicALU.png" alt="SchematicALU.png" />
<figcaption aria-hidden="true">SchematicALU.png</figcaption>
</figure>
<h3 id="instruction-memory">Instruction Memory</h3>
<p>The following is the schematic for the instruction memory of the MIPS
processor.</p>
<figure>
<img src="assets/SchematicInstructionMemory.png"
alt="SchematicInstructionMemory.png" />
<figcaption
aria-hidden="true">SchematicInstructionMemory.png</figcaption>
</figure>
<h3 id="program-counter">Program Counter</h3>
<p>The following is the schematic for the program counter of the MIPS
processor.</p>
<figure>
<img src="assets/SchematicProgramCounter.png"
alt="SchematicProgramCounter.png" />
<figcaption aria-hidden="true">SchematicProgramCounter.png</figcaption>
</figure>
<h3 id="waveform">Waveform</h3>
<p>The following is the wave form of the Processor from
modelsim/questasim:</p>
<figure>
<img src="assets/TestbenchWaveform.png" alt="TestbenchWaveform.png" />
<figcaption aria-hidden="true">TestbenchWaveform.png</figcaption>
</figure>
<h1 id="tooling">Tooling</h1>
<p>First, as learned in CPRE381, I enjoy having test-benches for my
code.</p>
<p>Test-benches allow for faster debugging and more efficient
development by allowing you to test your code without having to run it
on the FPGA board, directly seeing the signals, how they interact with
one another, and allows for testing to make sure that progress is being
made.</p>
<p>I used a test-bench to test my processor and ensure that it was
working correctly.</p>
<p>The following is the main test-bench that I used to test my
processor, <code>mips_tb.v</code>.</p>
<div class="sourceCode" id="cb70"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb70-1"><a href="#cb70-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb70-2"><a href="#cb70-2" aria-hidden="true" tabindex="-1"></a><span class="ot">`define CYCLE_TIME 20</span></span>
<span id="cb70-3"><a href="#cb70-3" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> mips_tb<span class="op">;</span></span>
<span id="cb70-4"><a href="#cb70-4" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> clk<span class="op">;</span></span>
<span id="cb70-5"><a href="#cb70-5" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> rst<span class="op">;</span></span>
<span id="cb70-6"><a href="#cb70-6" aria-hidden="true" tabindex="-1"></a>  <span class="dt">wire</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> seg_first<span class="op">,</span> seg_second<span class="op">,</span> seg_third<span class="op">,</span> seg_fourth<span class="op">,</span> seg_fifth<span class="op">;</span></span>
<span id="cb70-7"><a href="#cb70-7" aria-hidden="true" tabindex="-1"></a>  <span class="dt">integer</span> i<span class="op">;</span></span>
<span id="cb70-8"><a href="#cb70-8" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> #<span class="op">(</span><span class="ot">`CYCLE_TIME</span> <span class="op">/</span> <span class="dv">2</span><span class="op">)</span> clk <span class="op">=</span> <span class="op">~</span>clk<span class="op">;</span></span>
<span id="cb70-9"><a href="#cb70-9" aria-hidden="true" tabindex="-1"></a>  mips uut <span class="op">(</span></span>
<span id="cb70-10"><a href="#cb70-10" aria-hidden="true" tabindex="-1"></a>      .i_Clk<span class="op">(</span>clk<span class="op">),</span></span>
<span id="cb70-11"><a href="#cb70-11" aria-hidden="true" tabindex="-1"></a>      .i_Rst<span class="op">(</span>rst<span class="op">),</span></span>
<span id="cb70-12"><a href="#cb70-12" aria-hidden="true" tabindex="-1"></a>      .o_Seg_first<span class="op">(</span>seg_first<span class="op">),</span></span>
<span id="cb70-13"><a href="#cb70-13" aria-hidden="true" tabindex="-1"></a>      .o_Seg_second<span class="op">(</span>seg_second<span class="op">),</span></span>
<span id="cb70-14"><a href="#cb70-14" aria-hidden="true" tabindex="-1"></a>      .o_Seg_third<span class="op">(</span>seg_third<span class="op">),</span></span>
<span id="cb70-15"><a href="#cb70-15" aria-hidden="true" tabindex="-1"></a>      .o_Seg_fourth<span class="op">(</span>seg_fourth<span class="op">),</span></span>
<span id="cb70-16"><a href="#cb70-16" aria-hidden="true" tabindex="-1"></a>      .o_Seg_fifth<span class="op">(</span>seg_fifth<span class="op">)</span></span>
<span id="cb70-17"><a href="#cb70-17" aria-hidden="true" tabindex="-1"></a>  <span class="op">);</span></span>
<span id="cb70-18"><a href="#cb70-18" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb70-19"><a href="#cb70-19" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize data memory</span></span>
<span id="cb70-20"><a href="#cb70-20" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb70-21"><a href="#cb70-21" aria-hidden="true" tabindex="-1"></a>      uut.inst_DataMemory.Dmem<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb70-22"><a href="#cb70-22" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb70-23"><a href="#cb70-23" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize Register File</span></span>
<span id="cb70-24"><a href="#cb70-24" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb70-25"><a href="#cb70-25" aria-hidden="true" tabindex="-1"></a>      uut.inst_RegisterFile.RegData<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb70-26"><a href="#cb70-26" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb70-27"><a href="#cb70-27" aria-hidden="true" tabindex="-1"></a>    clk <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb70-28"><a href="#cb70-28" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb70-29"><a href="#cb70-29" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb70-30"><a href="#cb70-30" aria-hidden="true" tabindex="-1"></a>    <span class="bn">#1800</span> <span class="dt">$finish</span><span class="op">;</span></span>
<span id="cb70-31"><a href="#cb70-31" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb70-32"><a href="#cb70-32" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>I used do files to more easily compile and run my code.</p>
<p>The following is the main <code>.do</code> file, <code>run.do</code>
that I used to compile and run my code.</p>
<div class="sourceCode" id="cb71"><pre
class="sourceCode do"><code class="sourceCode stata"><span id="cb71-1"><a href="#cb71-1" aria-hidden="true" tabindex="-1"></a><span class="kw">set</span> target <span class="st">&quot;mips_tb&quot;</span></span>
<span id="cb71-2"><a href="#cb71-2" aria-hidden="true" tabindex="-1"></a><span class="kw">set</span> file <span class="st">&quot;proj/${target}.v&quot;</span></span>
<span id="cb71-3"><a href="#cb71-3" aria-hidden="true" tabindex="-1"></a><span class="kw">if</span> { [file exists <span class="st">&quot;work&quot;</span>] } {</span>
<span id="cb71-4"><a href="#cb71-4" aria-hidden="true" tabindex="-1"></a>    vdel -<span class="ot">all</span></span>
<span id="cb71-5"><a href="#cb71-5" aria-hidden="true" tabindex="-1"></a>}</span>
<span id="cb71-6"><a href="#cb71-6" aria-hidden="true" tabindex="-1"></a>vlog *.v</span>
<span id="cb71-7"><a href="#cb71-7" aria-hidden="true" tabindex="-1"></a>vsim -voptargs=+acc -debugDB <span class="ot">$target</span></span>
<span id="cb71-8"><a href="#cb71-8" aria-hidden="true" tabindex="-1"></a><span class="kw">force</span> -freeze sim:/<span class="ot">$target</span>/clk 1 0, 0 {5 ps} -<span class="fu">r</span> 10</span>
<span id="cb71-9"><a href="#cb71-9" aria-hidden="true" tabindex="-1"></a># <span class="kw">force</span> -freeze sim:/<span class="ot">$target</span>/rst 0 0, 1 {80 ps} -<span class="fu">r</span> 100</span>
<span id="cb71-10"><a href="#cb71-10" aria-hidden="true" tabindex="-1"></a>add wave -<span class="bn">position</span> insertpoint \ ../<span class="ot">$target</span><span class="co">/*</span></span>
<span id="cb71-11"><a href="#cb71-11" aria-hidden="true" tabindex="-1"></a><span class="co">run 1200</span></span></code></pre></div>
<p>In addition to the tooling already mentioned, I used
modelsim/questasim to simulate my processor and test-bench.</p>
<p>Furthermore, I used the Quartus Prime software to compile my code and
program my FPGA board.</p>
<p>Even further, I used tools like <a
href="https://git-scm.com/">git</a>, <a
href="https://github.com">GitHub</a>, and <a
href="https://www.markdownguide.org/">markdown</a> to manage my code,
version control, and documentation.</p>
<p>For editor tooling, I used <a href="https://neovim.io/">NeoVim</a>
with a combination of popular language servers that are used for
VerilogHDL.</p>
<p>These language servers that I used for development include <a
href="https://github.com/chipsalliance/verible">verible</a>, <a
href="https://tree-sitter.github.io/tree-sitter/">Tree-Sitter</a>, <a
href="https://github.com/vivekmalneedi/veridian">veridian</a>, and more
to provide completion, syntax highlighting, code actions, linting, and
more.</p>
<p>I think that using these language servers and tools in combination
with NeoVim (my personal open-sourced config has a startup time of
&lt;80ms) allowed me to develop my processor more efficiently and
effectively.</p>
<h2 id="components-and-explanations">Components and Explanations</h2>
<p>The following section explains the components of the MIPS processor
and their functionalities.</p>
<p>This is done by providing the Verilog code for each component and
explaining its role in the processor.</p>
<h4 id="alu-1">ALU</h4>
<p>The following is the code for the ALU module of the MIPS processor
called <code>ALU.v</code>. (It can be found in the <code>./proj/</code>
directory)</p>
<div class="sourceCode" id="cb72"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb72-1"><a href="#cb72-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb72-2"><a href="#cb72-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ALU <span class="op">(</span></span>
<span id="cb72-3"><a href="#cb72-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_data1<span class="op">,</span>        <span class="co">// data 1</span></span>
<span id="cb72-4"><a href="#cb72-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_read2<span class="op">,</span>        <span class="co">// data 2 from MUX</span></span>
<span id="cb72-5"><a href="#cb72-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Instruction<span class="op">,</span>  <span class="co">// used for sign-extension</span></span>
<span id="cb72-6"><a href="#cb72-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>             i_ALUSrc<span class="op">,</span></span>
<span id="cb72-7"><a href="#cb72-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span> <span class="dv">3</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUcontrol<span class="op">,</span></span>
<span id="cb72-8"><a href="#cb72-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span>        o_Zero<span class="op">,</span></span>
<span id="cb72-9"><a href="#cb72-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_ALUresult</span>
<span id="cb72-10"><a href="#cb72-10" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb72-11"><a href="#cb72-11" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> data2<span class="op">;</span></span>
<span id="cb72-12"><a href="#cb72-12" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_ALUSrc<span class="op">,</span> i_read2<span class="op">,</span> i_Instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb72-13"><a href="#cb72-13" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_ALUSrc <span class="op">==</span> <span class="dv">0</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb72-14"><a href="#cb72-14" aria-hidden="true" tabindex="-1"></a>      data2 <span class="op">=</span> i_read2<span class="op">;</span></span>
<span id="cb72-15"><a href="#cb72-15" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb72-16"><a href="#cb72-16" aria-hidden="true" tabindex="-1"></a>      <span class="kw">if</span> <span class="op">(</span>i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">]</span> <span class="op">==</span> <span class="bn">1&#39;b0</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb72-17"><a href="#cb72-17" aria-hidden="true" tabindex="-1"></a>        data2 <span class="op">=</span> <span class="op">{</span><span class="bn">16&#39;b0</span><span class="op">,</span> i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]};</span></span>
<span id="cb72-18"><a href="#cb72-18" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb72-19"><a href="#cb72-19" aria-hidden="true" tabindex="-1"></a>        data2 <span class="op">=</span> <span class="op">{{</span><span class="dv">16</span><span class="op">{</span><span class="bn">1&#39;b1</span><span class="op">}},</span> i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]};</span></span>
<span id="cb72-20"><a href="#cb72-20" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb72-21"><a href="#cb72-21" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb72-22"><a href="#cb72-22" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb72-23"><a href="#cb72-23" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> data2<span class="op">,</span> i_ALUcontrol<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb72-24"><a href="#cb72-24" aria-hidden="true" tabindex="-1"></a>    <span class="kw">case</span> <span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb72-25"><a href="#cb72-25" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0000</span><span class="op">:</span>  <span class="co">// AND</span></span>
<span id="cb72-26"><a href="#cb72-26" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">&amp;</span> data2<span class="op">;</span> <span class="co">// bitwise AND</span></span>
<span id="cb72-27"><a href="#cb72-27" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0001</span><span class="op">:</span>  <span class="co">// OR</span></span>
<span id="cb72-28"><a href="#cb72-28" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">|</span> data2<span class="op">;</span> <span class="co">// bitwise OR</span></span>
<span id="cb72-29"><a href="#cb72-29" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0010</span><span class="op">:</span>  <span class="co">// ADD</span></span>
<span id="cb72-30"><a href="#cb72-30" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">+</span> data2<span class="op">;</span> <span class="co">// addition</span></span>
<span id="cb72-31"><a href="#cb72-31" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0110</span><span class="op">:</span>  <span class="co">// SUB</span></span>
<span id="cb72-32"><a href="#cb72-32" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">-</span> data2<span class="op">;</span> <span class="co">// subtraction</span></span>
<span id="cb72-33"><a href="#cb72-33" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0111</span><span class="op">:</span>  <span class="co">// SLT</span></span>
<span id="cb72-34"><a href="#cb72-34" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> <span class="op">(</span>i_data1 <span class="op">&lt;</span> data2<span class="op">)</span> <span class="op">?</span> <span class="dv">1</span> <span class="op">:</span> <span class="dv">0</span><span class="op">;</span> <span class="co">// set-on-less-than</span></span>
<span id="cb72-35"><a href="#cb72-35" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b1100</span><span class="op">:</span>  <span class="co">// NOR</span></span>
<span id="cb72-36"><a href="#cb72-36" aria-hidden="true" tabindex="-1"></a>      o_ALUresult <span class="op">=</span> i_data1 <span class="op">|</span> <span class="op">~</span>data2<span class="op">;</span> <span class="co">// bitwise NOR</span></span>
<span id="cb72-37"><a href="#cb72-37" aria-hidden="true" tabindex="-1"></a>      <span class="kw">default</span><span class="op">:</span> <span class="op">;</span></span>
<span id="cb72-38"><a href="#cb72-38" aria-hidden="true" tabindex="-1"></a>    <span class="kw">endcase</span></span>
<span id="cb72-39"><a href="#cb72-39" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>o_ALUresult <span class="op">==</span> <span class="dv">0</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb72-40"><a href="#cb72-40" aria-hidden="true" tabindex="-1"></a>      o_Zero <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb72-41"><a href="#cb72-41" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb72-42"><a href="#cb72-42" aria-hidden="true" tabindex="-1"></a>      o_Zero <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb72-43"><a href="#cb72-43" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb72-44"><a href="#cb72-44" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb72-45"><a href="#cb72-45" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>The above code represents the ALU (Arithmetic Logic Unit) module of
the single-cycle MIPS processor.</p>
<p>The ALU is responsible for performing arithmetic and logical
operations on the input data based on the ALU control signal.</p>
<p>Inputs: - <code>i_data1</code> (32-bit): The first input data for the
ALU operation. - <code>i_read2</code> (32-bit): The second input data
from the MUX. - <code>i_Instruction</code> (32-bit): The instruction
used for sign-extension. - <code>i_ALUSrc</code> (1-bit): A control
signal indicating whether to use the second input data from the MUX or
the sign-extended immediate value. - <code>i_ALUcontrol</code> (4-bit):
The ALU control signal that determines the specific operation to be
performed.</p>
<p>Outputs: - <code>o_Zero</code> (1-bit): A flag indicating whether the
ALU result is zero. - <code>o_ALUresult</code> (32-bit): The result of
the ALU operation.</p>
<p>Functionality: 1. Data Selection: - The module first determines the
second input data for the ALU operation based on the
<code>i_ALUSrc</code> control signal. - If <code>i_ALUSrc</code> is 0,
the second input data is taken directly from <code>i_read2</code>. - If
<code>i_ALUSrc</code> is 1, the second input data is obtained by
sign-extending the 16-bit immediate value from the
<code>i_Instruction</code>.</p>
<ol start="2" type="1">
<li>ALU Operation:
<ul>
<li>Based on the <code>i_ALUcontrol</code> signal, the module performs
the corresponding ALU operation on <code>i_data1</code> and the selected
second input data (<code>data2</code>).</li>
<li>The supported ALU operations include AND, OR, ADD, SUB (subtract),
SLT (set-on-less-than), and NOR.</li>
<li>The result of the ALU operation is stored in
<code>o_ALUresult</code>.</li>
</ul></li>
<li>Zero Flag:
<ul>
<li>After performing the ALU operation, the module checks if the result
is zero.</li>
<li>If the result is zero, the <code>o_Zero</code> flag is set to 1;
otherwise, it is set to 0.</li>
</ul></li>
</ol>
<p>Significance in the MIPS Processor: - The ALU is a crucial component
in the MIPS processor’s data-path. - It performs arithmetic and logical
operations on the input data based on the instructions being executed. -
The ALU receives input data from the register file or the immediate
value in the instruction, depending on the <code>i_ALUSrc</code> control
signal. - The ALU control signal (<code>i_ALUcontrol</code>) determines
the specific operation to be performed, which is decoded by the ALU
control module based on the instruction opcode and function code. - The
result of the ALU operation is used for various purposes, such as
storing it back to the register file, using it as a memory address, or
making branching decisions based on the zero flag.</p>
<p>Interaction with Other Components: - The ALU receives input data from
the register file (<code>i_data1</code> and <code>i_read2</code>) and
the instruction (<code>i_Instruction</code>). - The ALU control module
generates the <code>i_ALUcontrol</code> signal based on the instruction
opcode and function code, which determines the ALU operation to be
performed. - The ALU result (<code>o_ALUresult</code>) is used by other
components, such as the data memory for memory operations or the
register file for storing the result. - The zero flag
(<code>o_Zero</code>) is used by the control unit to make branching
decisions based on the comparison result.</p>
<p>Overall, the ALU module performs the necessary arithmetic and logical
operations in the MIPS processor based on the instruction being
executed.</p>
<p>The ALU allows the processor to execute instructions and produce the
desired results.</p>
<h4 id="control-unit">Control Unit</h4>
<p>The following is the code for the control unit of the MIPS processor
called <code>ControlUnit.v</code>. (It can be found in the
<code>./proj/</code> directory)</p>
<p>As named, the <code>ControlUnit</code>, <code>ControlUnit.v</code> is
responsible for decoding the instruction and generating the control
signals for the various components of the processor.</p>
<div class="sourceCode" id="cb73"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb73-1"><a href="#cb73-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb73-2"><a href="#cb73-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ControlUnit <span class="op">(</span></span>
<span id="cb73-3"><a href="#cb73-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_instruction<span class="op">,</span></span>
<span id="cb73-4"><a href="#cb73-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_RegDst<span class="op">,</span></span>
<span id="cb73-5"><a href="#cb73-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_Jump<span class="op">,</span></span>
<span id="cb73-6"><a href="#cb73-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_Branch<span class="op">,</span></span>
<span id="cb73-7"><a href="#cb73-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_Bne<span class="op">,</span></span>
<span id="cb73-8"><a href="#cb73-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_MemRead<span class="op">,</span></span>
<span id="cb73-9"><a href="#cb73-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_MemtoReg<span class="op">,</span></span>
<span id="cb73-10"><a href="#cb73-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">1</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_ALUOp<span class="op">,</span></span>
<span id="cb73-11"><a href="#cb73-11" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_MemWrite<span class="op">,</span></span>
<span id="cb73-12"><a href="#cb73-12" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_ALUSrc<span class="op">,</span></span>
<span id="cb73-13"><a href="#cb73-13" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_RegWrite<span class="op">,</span></span>
<span id="cb73-14"><a href="#cb73-14" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_first<span class="op">,</span></span>
<span id="cb73-15"><a href="#cb73-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_second<span class="op">,</span></span>
<span id="cb73-16"><a href="#cb73-16" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_third<span class="op">,</span></span>
<span id="cb73-17"><a href="#cb73-17" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_fourth<span class="op">,</span></span>
<span id="cb73-18"><a href="#cb73-18" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_fifth</span>
<span id="cb73-19"><a href="#cb73-19" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb73-20"><a href="#cb73-20" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb73-21"><a href="#cb73-21" aria-hidden="true" tabindex="-1"></a>    o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-22"><a href="#cb73-22" aria-hidden="true" tabindex="-1"></a>    o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-23"><a href="#cb73-23" aria-hidden="true" tabindex="-1"></a>    o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-24"><a href="#cb73-24" aria-hidden="true" tabindex="-1"></a>    o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-25"><a href="#cb73-25" aria-hidden="true" tabindex="-1"></a>    o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-26"><a href="#cb73-26" aria-hidden="true" tabindex="-1"></a>    o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb73-27"><a href="#cb73-27" aria-hidden="true" tabindex="-1"></a>    o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-28"><a href="#cb73-28" aria-hidden="true" tabindex="-1"></a>    o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-29"><a href="#cb73-29" aria-hidden="true" tabindex="-1"></a>    o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-30"><a href="#cb73-30" aria-hidden="true" tabindex="-1"></a>    o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-31"><a href="#cb73-31" aria-hidden="true" tabindex="-1"></a>    o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-32"><a href="#cb73-32" aria-hidden="true" tabindex="-1"></a>    o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-33"><a href="#cb73-33" aria-hidden="true" tabindex="-1"></a>    o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-34"><a href="#cb73-34" aria-hidden="true" tabindex="-1"></a>    o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-35"><a href="#cb73-35" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb73-36"><a href="#cb73-36" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb73-37"><a href="#cb73-37" aria-hidden="true" tabindex="-1"></a>    <span class="kw">case</span> <span class="op">(</span>i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">])</span></span>
<span id="cb73-38"><a href="#cb73-38" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000000</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// ARITHMETIC</span></span>
<span id="cb73-39"><a href="#cb73-39" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-40"><a href="#cb73-40" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-41"><a href="#cb73-41" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-42"><a href="#cb73-42" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-43"><a href="#cb73-43" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-44"><a href="#cb73-44" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-45"><a href="#cb73-45" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-46"><a href="#cb73-46" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-47"><a href="#cb73-47" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b10</span><span class="op">;</span></span>
<span id="cb73-48"><a href="#cb73-48" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-49"><a href="#cb73-49" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span>  <span class="bn">7&#39;b0001000</span><span class="op">;</span>  <span class="co">// A</span></span>
<span id="cb73-50"><a href="#cb73-50" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111010</span><span class="op">;</span>  <span class="co">// R</span></span>
<span id="cb73-51"><a href="#cb73-51" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span>  <span class="bn">7&#39;b1111001</span><span class="op">;</span>  <span class="co">// I</span></span>
<span id="cb73-52"><a href="#cb73-52" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b0001111</span><span class="op">;</span>  <span class="co">// T</span></span>
<span id="cb73-53"><a href="#cb73-53" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span>  <span class="bn">7&#39;b0001001</span><span class="op">;</span>  <span class="co">// H</span></span>
<span id="cb73-54"><a href="#cb73-54" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-55"><a href="#cb73-55" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b001000</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// addi</span></span>
<span id="cb73-56"><a href="#cb73-56" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-57"><a href="#cb73-57" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-58"><a href="#cb73-58" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-59"><a href="#cb73-59" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-60"><a href="#cb73-60" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-61"><a href="#cb73-61" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-62"><a href="#cb73-62" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-63"><a href="#cb73-63" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-64"><a href="#cb73-64" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb73-65"><a href="#cb73-65" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-66"><a href="#cb73-66" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b0001000</span><span class="op">;</span>  <span class="co">// A</span></span>
<span id="cb73-67"><a href="#cb73-67" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1000010</span><span class="op">;</span>  <span class="co">// d</span></span>
<span id="cb73-68"><a href="#cb73-68" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1000010</span><span class="op">;</span>  <span class="co">// d</span></span>
<span id="cb73-69"><a href="#cb73-69" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1001111</span><span class="op">;</span>  <span class="co">// i</span></span>
<span id="cb73-70"><a href="#cb73-70" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-71"><a href="#cb73-71" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-72"><a href="#cb73-72" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b001100</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// andi</span></span>
<span id="cb73-73"><a href="#cb73-73" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-74"><a href="#cb73-74" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-75"><a href="#cb73-75" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-76"><a href="#cb73-76" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-77"><a href="#cb73-77" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-78"><a href="#cb73-78" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-79"><a href="#cb73-79" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-80"><a href="#cb73-80" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-81"><a href="#cb73-81" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b11</span><span class="op">;</span></span>
<span id="cb73-82"><a href="#cb73-82" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-83"><a href="#cb73-83" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b0001000</span><span class="op">;</span>  <span class="co">// A</span></span>
<span id="cb73-84"><a href="#cb73-84" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b0101011</span><span class="op">;</span>  <span class="co">// n</span></span>
<span id="cb73-85"><a href="#cb73-85" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1000010</span><span class="op">;</span>  <span class="co">// d</span></span>
<span id="cb73-86"><a href="#cb73-86" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1001111</span><span class="op">;</span>  <span class="co">// i</span></span>
<span id="cb73-87"><a href="#cb73-87" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-88"><a href="#cb73-88" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-89"><a href="#cb73-89" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b100011</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// lw</span></span>
<span id="cb73-90"><a href="#cb73-90" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-91"><a href="#cb73-91" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-92"><a href="#cb73-92" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-93"><a href="#cb73-93" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-94"><a href="#cb73-94" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-95"><a href="#cb73-95" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-96"><a href="#cb73-96" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-97"><a href="#cb73-97" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-98"><a href="#cb73-98" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb73-99"><a href="#cb73-99" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-100"><a href="#cb73-100" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1000111</span><span class="op">;</span>  <span class="co">// L</span></span>
<span id="cb73-101"><a href="#cb73-101" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1001001</span><span class="op">;</span>  <span class="co">// w</span></span>
<span id="cb73-102"><a href="#cb73-102" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-103"><a href="#cb73-103" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-104"><a href="#cb73-104" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-105"><a href="#cb73-105" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-106"><a href="#cb73-106" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b101011</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// sw</span></span>
<span id="cb73-107"><a href="#cb73-107" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-108"><a href="#cb73-108" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-109"><a href="#cb73-109" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-110"><a href="#cb73-110" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-111"><a href="#cb73-111" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-112"><a href="#cb73-112" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-113"><a href="#cb73-113" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-114"><a href="#cb73-114" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-115"><a href="#cb73-115" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb73-116"><a href="#cb73-116" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-117"><a href="#cb73-117" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b0010010</span><span class="op">;</span>  <span class="co">// S</span></span>
<span id="cb73-118"><a href="#cb73-118" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1001001</span><span class="op">;</span>  <span class="co">// w</span></span>
<span id="cb73-119"><a href="#cb73-119" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-120"><a href="#cb73-120" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-121"><a href="#cb73-121" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-122"><a href="#cb73-122" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-123"><a href="#cb73-123" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000100</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// beq</span></span>
<span id="cb73-124"><a href="#cb73-124" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-125"><a href="#cb73-125" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-126"><a href="#cb73-126" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-127"><a href="#cb73-127" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-128"><a href="#cb73-128" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-129"><a href="#cb73-129" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-130"><a href="#cb73-130" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-131"><a href="#cb73-131" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-132"><a href="#cb73-132" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b01</span><span class="op">;</span></span>
<span id="cb73-133"><a href="#cb73-133" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-134"><a href="#cb73-134" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1100000</span><span class="op">;</span>  <span class="co">// b</span></span>
<span id="cb73-135"><a href="#cb73-135" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b0110000</span><span class="op">;</span>  <span class="co">// e</span></span>
<span id="cb73-136"><a href="#cb73-136" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b0001100</span><span class="op">;</span>  <span class="co">// q</span></span>
<span id="cb73-137"><a href="#cb73-137" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-138"><a href="#cb73-138" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-139"><a href="#cb73-139" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-140"><a href="#cb73-140" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000101</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// bne</span></span>
<span id="cb73-141"><a href="#cb73-141" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-142"><a href="#cb73-142" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-143"><a href="#cb73-143" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-144"><a href="#cb73-144" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-145"><a href="#cb73-145" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-146"><a href="#cb73-146" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-147"><a href="#cb73-147" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-148"><a href="#cb73-148" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-149"><a href="#cb73-149" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b01</span><span class="op">;</span></span>
<span id="cb73-150"><a href="#cb73-150" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-151"><a href="#cb73-151" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1100000</span><span class="op">;</span>  <span class="co">// b</span></span>
<span id="cb73-152"><a href="#cb73-152" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b0101011</span><span class="op">;</span>  <span class="co">// n</span></span>
<span id="cb73-153"><a href="#cb73-153" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b0110000</span><span class="op">;</span>  <span class="co">// e</span></span>
<span id="cb73-154"><a href="#cb73-154" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-155"><a href="#cb73-155" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-156"><a href="#cb73-156" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-157"><a href="#cb73-157" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000010</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// j</span></span>
<span id="cb73-158"><a href="#cb73-158" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-159"><a href="#cb73-159" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-160"><a href="#cb73-160" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-161"><a href="#cb73-161" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-162"><a href="#cb73-162" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-163"><a href="#cb73-163" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-164"><a href="#cb73-164" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-165"><a href="#cb73-165" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-166"><a href="#cb73-166" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b01</span><span class="op">;</span></span>
<span id="cb73-167"><a href="#cb73-167" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb73-168"><a href="#cb73-168" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1100001</span><span class="op">;</span>  <span class="co">// J</span></span>
<span id="cb73-169"><a href="#cb73-169" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-170"><a href="#cb73-170" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-171"><a href="#cb73-171" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-172"><a href="#cb73-172" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-173"><a href="#cb73-173" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-174"><a href="#cb73-174" aria-hidden="true" tabindex="-1"></a>      <span class="kw">default</span><span class="op">:</span> <span class="kw">begin</span></span>
<span id="cb73-175"><a href="#cb73-175" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-176"><a href="#cb73-176" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-177"><a href="#cb73-177" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb73-178"><a href="#cb73-178" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-179"><a href="#cb73-179" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-180"><a href="#cb73-180" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-181"><a href="#cb73-181" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-182"><a href="#cb73-182" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-183"><a href="#cb73-183" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb73-184"><a href="#cb73-184" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb73-185"><a href="#cb73-185" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-186"><a href="#cb73-186" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-187"><a href="#cb73-187" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-188"><a href="#cb73-188" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-189"><a href="#cb73-189" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb73-190"><a href="#cb73-190" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb73-191"><a href="#cb73-191" aria-hidden="true" tabindex="-1"></a>    <span class="kw">endcase</span></span>
<span id="cb73-192"><a href="#cb73-192" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb73-193"><a href="#cb73-193" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>The above Verilog code represents the Control Unit component of the
single-cycle MIPS processor.</p>
<p>The Control Unit is responsible for generating control signals based
on the input instruction, which determine the behavior of various
components within the processor.</p>
<h5 id="io">IO</h5>
<p>The Control Unit, <code>ControlUnit.v</code>, has the following
inputs and outputs:</p>
<p>Inputs: - <code>i_instruction</code>: The 32-bit instruction fetched
from the Instruction Memory.</p>
<p>Outputs: - Various control signals: - <code>o_RegDst</code>: Selects
the destination register for the instruction (0 for rt, 1 for rd). -
<code>o_Jump</code>: Indicates if the instruction is a jump instruction.
- <code>o_Branch</code>: Indicates if the instruction is a branch
instruction. - <code>o_Bne</code>: Indicates if the instruction is a
“branch not equal” instruction. - <code>o_MemRead</code>: Enables
reading from the Data Memory. - <code>o_MemtoReg</code>: Selects the
source of data to be written to the register (0 for ALU result, 1 for
memory data). - <code>o_ALUOp</code>: A 2-bit signal that specifies the
ALU operation. - <code>o_MemWrite</code>: Enables writing to the Data
Memory. - <code>o_ALUSrc</code>: Selects the second source for the ALU
(0 for register, 1 for immediate). - <code>o_RegWrite</code>: Enables
writing to the Register File. - 7-segment display outputs: -
<code>o_seg_first</code> to <code>o_seg_fifth</code>: Control signals
for displaying the instruction type on 7-segment displays.</p>
<h5 id="processor-context">Processor Context:</h5>
<p>The following is the context (the purpose and functionality) of the
Control Unit in the single-cycle MIPS processor:</p>
<ol type="1">
<li><p>The Control Unit initializes all control signals to default
values in the <code>initial</code> block.</p></li>
<li><p>The <code>always</code> block is triggered whenever the
<code>i_instruction</code> changes. It uses a case statement to
determine the type of instruction based on the opcode (bits 31 to 26 of
the instruction).</p></li>
<li><p>Depending on the instruction type, the Control Unit sets the
appropriate control signals:</p>
<ul>
<li>For R-type instructions (arithmetic), it sets <code>RegDst</code> to
1, enables <code>RegWrite</code>, sets <code>ALUOp</code> to 2’b10, and
displays “ARITH” on the 7-segment displays.</li>
<li>For I-type instructions (addi, andi, lw, sw), it sets
<code>ALUSrc</code> to 1, enables <code>RegWrite</code> (except for sw),
sets <code>ALUOp</code> based on the instruction, and displays the
instruction type on the 7-segment displays.</li>
<li>For branch instructions (beq, bne), it sets <code>Branch</code> to
1, sets <code>ALUOp</code> to 2’b01, and displays the instruction type
on the 7-segment displays.</li>
<li>For the jump instruction, it sets <code>Jump</code> to 1, sets
<code>ALUOp</code> to 2’b01, and displays “J” on the 7-segment
displays.</li>
</ul></li>
<li><p>If the instruction does not match any of the defined cases, the
Control Unit sets all control signals to their default values and
displays blank on the 7-segment displays.</p></li>
</ol>
<p>The Control Unit is critial to correctly orchestrating the operation
of the single-cycle MIPS processor.</p>
<p>It interprets the instruction and generates the necessary control
signals to control the data-path components, such as the ALU, Register
File, and Data Memory.</p>
<p>The control signals determine the flow of data and the operations
performed in each stage of the processor pipeline.</p>
<ul>
<li>The Control Unit receives the instruction from the Instruction
Memory.</li>
<li>It sends control signals to various components, such as the ALU,
Register File, and Data Memory, to control their behavior based on the
instruction being executed.</li>
<li>The control signals generated by the Control Unit are used by the
data-path components to perform the required operations and route the
data accordingly.</li>
</ul>
<h4 id="data-memory-1">Data Memory</h4>
<p>The following is the code for the Data Memory module in the MIPS
processor called <code>DataMemory.v</code>. (It can be found in the
<code>./proj/</code> directory)</p>
<div class="sourceCode" id="cb74"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb74-1"><a href="#cb74-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb74-2"><a href="#cb74-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> DataMemory <span class="op">(</span></span>
<span id="cb74-3"><a href="#cb74-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_clk<span class="op">,</span></span>
<span id="cb74-4"><a href="#cb74-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_addr<span class="op">,</span></span>
<span id="cb74-5"><a href="#cb74-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_wData<span class="op">,</span></span>
<span id="cb74-6"><a href="#cb74-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUresult<span class="op">,</span></span>
<span id="cb74-7"><a href="#cb74-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemWrite<span class="op">,</span></span>
<span id="cb74-8"><a href="#cb74-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemRead<span class="op">,</span></span>
<span id="cb74-9"><a href="#cb74-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemtoReg<span class="op">,</span></span>
<span id="cb74-10"><a href="#cb74-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_rData</span>
<span id="cb74-11"><a href="#cb74-11" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb74-12"><a href="#cb74-12" aria-hidden="true" tabindex="-1"></a>  <span class="dt">parameter</span> SIZE_DM <span class="op">=</span> <span class="dv">128</span><span class="op">;</span>       <span class="co">// size of this memory, by default 128*32</span></span>
<span id="cb74-13"><a href="#cb74-13" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> Dmem<span class="op">[</span>SIZE_DM<span class="dv">-1</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span>  <span class="co">// instruction memory</span></span>
<span id="cb74-14"><a href="#cb74-14" aria-hidden="true" tabindex="-1"></a>  <span class="dt">integer</span> i<span class="op">;</span></span>
<span id="cb74-15"><a href="#cb74-15" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb74-16"><a href="#cb74-16" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> SIZE_DM<span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb74-17"><a href="#cb74-17" aria-hidden="true" tabindex="-1"></a>      Dmem<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb74-18"><a href="#cb74-18" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb74-19"><a href="#cb74-19" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb74-20"><a href="#cb74-20" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_addr <span class="dt">or</span> i_MemRead <span class="dt">or</span> i_MemtoReg <span class="dt">or</span> i_ALUresult<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb74-21"><a href="#cb74-21" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_MemRead <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb74-22"><a href="#cb74-22" aria-hidden="true" tabindex="-1"></a>      <span class="kw">if</span> <span class="op">(</span>i_MemtoReg <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb74-23"><a href="#cb74-23" aria-hidden="true" tabindex="-1"></a>        o_rData <span class="op">=</span> Dmem<span class="op">[</span>i_addr<span class="op">];</span></span>
<span id="cb74-24"><a href="#cb74-24" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb74-25"><a href="#cb74-25" aria-hidden="true" tabindex="-1"></a>        o_rData <span class="op">=</span> i_ALUresult<span class="op">;</span></span>
<span id="cb74-26"><a href="#cb74-26" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb74-27"><a href="#cb74-27" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb74-28"><a href="#cb74-28" aria-hidden="true" tabindex="-1"></a>      o_rData <span class="op">=</span> i_ALUresult<span class="op">;</span></span>
<span id="cb74-29"><a href="#cb74-29" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb74-30"><a href="#cb74-30" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb74-31"><a href="#cb74-31" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_clk<span class="op">)</span> <span class="kw">begin</span>  <span class="co">// MemWrite, wData, addr</span></span>
<span id="cb74-32"><a href="#cb74-32" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_MemWrite <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb74-33"><a href="#cb74-33" aria-hidden="true" tabindex="-1"></a>      Dmem<span class="op">[</span>i_addr<span class="op">]</span> <span class="op">=</span> i_wData<span class="op">;</span></span>
<span id="cb74-34"><a href="#cb74-34" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb74-35"><a href="#cb74-35" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb74-36"><a href="#cb74-36" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>The provided code snippet is the implementation of the Data Memory
module (<code>DataMemory.v</code>) in the single-cycle MIPS
processor.</p>
<p>The Data Memory module serves as the main memory for storing and
retrieving data in the processor.</p>
<h5 id="io-1">IO</h5>
<p>The following are the detailed input and output ports of the Data
Memory module, <code>DataMemory.v</code>: - <code>i_clk</code>: Input
clock signal. - <code>i_addr</code>: Input address for reading or
writing data. - <code>i_wData</code>: Input write data to be stored in
memory. - <code>i_ALUresult</code>: Input ALU result, which can be used
as the address or data depending on the control signals. -
<code>i_MemWrite</code>: Input control signal indicating a memory write
operation. - <code>i_MemRead</code>: Input control signal indicating a
memory read operation. - <code>i_MemtoReg</code>: Input control signal
indicating whether to pass the memory read data or ALU result to the
output. - <code>o_rData</code>: Output read data from the memory.</p>
<h5 id="functionality">Functionality</h5>
<p>Memory Initialization:</p>
<ul>
<li>The module defines a parameter <code>SIZE_DM</code> representing the
size of the data memory (default is 128 words).</li>
<li>It declares a register array <code>Dmem</code> of size
<code>SIZE_DM</code> to store the memory contents.</li>
<li>In the initial block, all memory locations are initialized to zero
using a loop.</li>
</ul>
<p>Memory Read Operation:</p>
<ul>
<li>The first always block is triggered whenever the input signals
<code>i_addr</code>, <code>i_MemRead</code>, <code>i_MemtoReg</code>, or
<code>i_ALUresult</code> change.</li>
<li>If <code>i_MemRead</code> is asserted (equals 1), it indicates a
memory read operation.
<ul>
<li>If <code>i_MemtoReg</code> is also asserted, the data at memory
location <code>i_addr</code> is assigned to the output
<code>o_rData</code>.</li>
<li>Otherwise, the ALU result <code>i_ALUresult</code> is assigned to
<code>o_rData</code>.</li>
</ul></li>
<li>If <code>i_MemRead</code> is not asserted, the ALU result
<code>i_ALUresult</code> is directly assigned to
<code>o_rData</code>.</li>
</ul>
<p>Memory Write Operation:</p>
<ul>
<li>The second always block is triggered on the positive edge of the
clock signal <code>i_clk</code>.</li>
<li>If <code>i_MemWrite</code> is asserted (equals 1), it indicates a
memory write operation.</li>
<li>The data <code>i_wData</code> is written to the memory location
specified by <code>i_addr</code>.</li>
</ul>
<h5 id="significance-in-the-processor">Significance in the
Processor</h5>
<p>Interaction with Other Components: - The Data Memory module interacts
with the ALU and the Control Unit in the processor. - The ALU provides
the address (<code>i_ALUresult</code>) for memory read or write
operations. - The Control Unit generates the control signals
(<code>i_MemWrite</code>, <code>i_MemRead</code>,
<code>i_MemtoReg</code>) to control the behavior of the Data Memory
module. - The Register File provides the data to be written to memory
(<code>i_wData</code>) during a memory write operation. - The output
read data (<code>o_rData</code>) is passed back to the Register File or
used as needed in subsequent stages of the processor pipeline.</p>
<p>Simply put, the Data Memory module allows storing and retrieving data
in the MIPS processor.</p>
<p>It responds to memory read and write requests based on the provided
address and control signals, and it interacts with other components such
as the ALU, Control Unit, and Register File to facilitate data storage
and retrieval operations.</p>
<h4 id="instruction-memory-1">Instruction Memory</h4>
<p>The following is the code for the Instruction Memory module in the
MIPS processor called <code>InstructionMemory.v</code>. (It can be found
in the <code>./proj/</code> directory)</p>
<div class="sourceCode" id="cb75"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb75-1"><a href="#cb75-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb75-2"><a href="#cb75-2" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> InstructionMemory <span class="op">(</span></span>
<span id="cb75-3"><a href="#cb75-3" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Addr<span class="op">,</span></span>
<span id="cb75-4"><a href="#cb75-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Ctr<span class="op">,</span>          <span class="co">// [31-26]</span></span>
<span id="cb75-5"><a href="#cb75-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Funcode<span class="op">,</span>      <span class="co">// [5-0]</span></span>
<span id="cb75-6"><a href="#cb75-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Instruction  <span class="co">// [31-0]</span></span>
<span id="cb75-7"><a href="#cb75-7" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb75-8"><a href="#cb75-8" aria-hidden="true" tabindex="-1"></a>  <span class="dt">parameter</span> SIZE_IM <span class="op">=</span> <span class="dv">128</span><span class="op">;</span>           <span class="co">// size of this memory, by default 128*32</span></span>
<span id="cb75-9"><a href="#cb75-9" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> Imem<span class="op">[</span>SIZE_IM<span class="dv">-1</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span>      <span class="co">// instruction memory</span></span>
<span id="cb75-10"><a href="#cb75-10" aria-hidden="true" tabindex="-1"></a>  <span class="dt">integer</span> n<span class="op">;</span></span>
<span id="cb75-11"><a href="#cb75-11" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb75-12"><a href="#cb75-12" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>n <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> n <span class="op">&lt;</span> SIZE_IM<span class="op">;</span> n <span class="op">=</span> n <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb75-13"><a href="#cb75-13" aria-hidden="true" tabindex="-1"></a>      Imem<span class="op">[</span>n<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b11111100000000000000000000000000</span><span class="op">;</span></span>
<span id="cb75-14"><a href="#cb75-14" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb75-15"><a href="#cb75-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">$readmemb</span><span class="op">(</span><span class="st">&quot;instructions.mem&quot;</span><span class="op">,</span> Imem<span class="op">);</span></span>
<span id="cb75-16"><a href="#cb75-16" aria-hidden="true" tabindex="-1"></a>    i_Instruction <span class="op">=</span> <span class="bn">32&#39;b11111100000000000000000000000000</span><span class="op">;</span></span>
<span id="cb75-17"><a href="#cb75-17" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb75-18"><a href="#cb75-18" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_Addr<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb75-19"><a href="#cb75-19" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_Addr <span class="op">==</span> <span class="op">-</span><span class="dv">4</span><span class="op">)</span> <span class="kw">begin</span>         <span class="co">// init</span></span>
<span id="cb75-20"><a href="#cb75-20" aria-hidden="true" tabindex="-1"></a>      i_Instruction <span class="op">=</span> <span class="bn">32&#39;b11111100000000000000000000000000</span><span class="op">;</span></span>
<span id="cb75-21"><a href="#cb75-21" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb75-22"><a href="#cb75-22" aria-hidden="true" tabindex="-1"></a>      i_Instruction <span class="op">=</span> Imem<span class="op">[</span>i_Addr<span class="op">&gt;&gt;</span><span class="dv">2</span><span class="op">];</span></span>
<span id="cb75-23"><a href="#cb75-23" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb75-24"><a href="#cb75-24" aria-hidden="true" tabindex="-1"></a>    i_Ctr <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">];</span></span>
<span id="cb75-25"><a href="#cb75-25" aria-hidden="true" tabindex="-1"></a>    i_Funcode <span class="op">=</span> i_Instruction<span class="op">[</span><span class="dv">5</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span></span>
<span id="cb75-26"><a href="#cb75-26" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb75-27"><a href="#cb75-27" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>The provided code represents “the Instruction Memory module
<code>InstructionMemory.v</code> for the single-cycle MIPS
processor.”</p>
<h5 id="purpose">Purpose:</h5>
<p>The Instruction Memory module is responsible for storing the
processor’s instructions and providing them to the other components of
the processor.</p>
<p>It acts as a read-only memory (ROM) that holds the program
instructions.</p>
<h5 id="io-2">IO</h5>
<p>Inputs and Outputs:</p>
<ul>
<li><code>i_Addr</code> (input, 32-bit): Represents the memory address
from which the instruction should be fetched.</li>
<li><code>i_Ctr</code> (output, 6-bit): Outputs the control bits of the
fetched instruction (bits [31:26]).</li>
<li><code>i_Funcode</code> (output, 6-bit): Outputs the function code of
the fetched instruction (bits [5:0]).</li>
<li><code>i_Instruction</code> (output, 32-bit): Outputs the complete
32-bit fetched instruction.</li>
</ul>
<h5 id="functionality-1">Functionality</h5>
<ol type="1">
<li>The module defines a parameter <code>SIZE_IM</code> which represents
the size of the instruction memory.</li>
</ol>
<p>By default, it is set to 128, meaning the memory can hold 128 32-bit
instructions.</p>
<ol start="2" type="1">
<li><p>The module declares a register array <code>Imem</code> of size
<code>SIZE_IM</code> to store the instructions.</p></li>
<li><p>In the initial block:</p>
<ul>
<li>The memory is initialized with a default instruction
(32’b11111100000000000000000000000000) using a loop.</li>
<li>The instructions are then loaded from a file named
“instructions.mem” using the <code>$readmemb</code> system task. This
file contains the binary representation of the instructions.</li>
<li>The <code>i_Instruction</code> output is initialized with the
default instruction.</li>
</ul></li>
<li><p>The module has an “always” block that is triggered whenever the
<code>i_Addr</code> input changes:</p>
<ul>
<li>If <code>i_Addr</code> is equal to -4 (used for initialization), the
<code>i_Instruction</code> output is set to the default
instruction.</li>
<li>Otherwise, the instruction is fetched from the <code>Imem</code>
array using the address <code>i_Addr</code> shifted right by 2 bits
(assuming word-aligned addresses).</li>
<li>The control bits (<code>i_Ctr</code>) and function code
(<code>i_Funcode</code>) are extracted from the fetched instruction and
assigned to the respective outputs.</li>
</ul></li>
</ol>
<h5 id="processor-context-1">Processor Context</h5>
<p>The following is the context of the Instruction Memory module,
<code>InstructionMemory.v</code> in the single-cycle MIPS processor,
<code>mips.v</code>:</p>
<ul>
<li>The Program Counter (<code>PC</code>) module provides the memory
address (<code>i_Addr</code>) to the Instruction Memory module to fetch
the instruction at that address.</li>
<li>The fetched instruction (<code>i_Instruction</code>) is then passed
to other components of the processor, such as the Control Unit and the
Register File, for further processing and execution.</li>
<li>The control bits (<code>i_Ctr</code>) and function code
(<code>i_Funcode</code>) are used by the Control Unit to generate
appropriate control signals for the processor’s data-path.</li>
</ul>
<h5 id="significance">Significance</h5>
<p>The Instruction Memory module is a crucial component of the MIPS
processor as it holds the program instructions that the processor
executes.</p>
<p>It provides the instructions to the processor’s data-path, enabling
the processor to perform the desired operations and execute the program
stored in the memory.</p>
<h5 id="summary">Summary</h5>
<p>Essentially, the Instruction Memory module in the single-cycle MIPS
processor acts as a read-only memory that stores the program
instructions. It fetches instructions based on the provided memory
address and outputs the complete instruction along with its control bits
and function code for further processing by other components of the
processor.</p>
<h2 id="verbose-components-code">Verbose Components Code</h2>
<p>The following section shows detailed, commented code files for each
of the components of the processor.</p>
<p>It includes detailed code comments to better explain the
functionality and purpose of each component within the actual verilog
code.</p>
<h3 id="data-memory-2">Data Memory</h3>
<p>The following is the commented code for the Data Memory module in the
MIPS processor called <code>DataMemory.v</code>:</p>
<div class="sourceCode" id="cb76"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb76-1"><a href="#cb76-1" aria-hidden="true" tabindex="-1"></a><span class="co">// File: DataMemory.v</span></span>
<span id="cb76-2"><a href="#cb76-2" aria-hidden="true" tabindex="-1"></a><span class="co">// Description: This file contains the data memory module for the MIPS processor.</span></span>
<span id="cb76-3"><a href="#cb76-3" aria-hidden="true" tabindex="-1"></a><span class="co">// Purpose: The data memory stores data values and provides read and write access to the processor.</span></span>
<span id="cb76-4"><a href="#cb76-4" aria-hidden="true" tabindex="-1"></a><span class="co">//          It is responsible for handling memory read and write operations based on the control signals</span></span>
<span id="cb76-5"><a href="#cb76-5" aria-hidden="true" tabindex="-1"></a><span class="co">//          received from the control unit.</span></span>
<span id="cb76-6"><a href="#cb76-6" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb76-7"><a href="#cb76-7" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> DataMemory <span class="op">(</span></span>
<span id="cb76-8"><a href="#cb76-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_clk<span class="op">,</span>                    <span class="co">// Clock input</span></span>
<span id="cb76-9"><a href="#cb76-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_addr<span class="op">,</span>            <span class="co">// Address input for memory access</span></span>
<span id="cb76-10"><a href="#cb76-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_wData<span class="op">,</span>           <span class="co">// Write data input</span></span>
<span id="cb76-11"><a href="#cb76-11" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUresult<span class="op">,</span>       <span class="co">// ALU result input (used for memory address calculation)</span></span>
<span id="cb76-12"><a href="#cb76-12" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemWrite<span class="op">,</span>               <span class="co">// Control signal for memory write operation</span></span>
<span id="cb76-13"><a href="#cb76-13" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemRead<span class="op">,</span>                <span class="co">// Control signal for memory read operation</span></span>
<span id="cb76-14"><a href="#cb76-14" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_MemtoReg<span class="op">,</span>               <span class="co">// Control signal for selecting memory or ALU result as the output</span></span>
<span id="cb76-15"><a href="#cb76-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_rData       <span class="co">// Read data output</span></span>
<span id="cb76-16"><a href="#cb76-16" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb76-17"><a href="#cb76-17" aria-hidden="true" tabindex="-1"></a>  <span class="dt">parameter</span> SIZE_DM <span class="op">=</span> <span class="dv">128</span><span class="op">;</span>           <span class="co">// Size of the data memory (default: 128 * 32 bits)</span></span>
<span id="cb76-18"><a href="#cb76-18" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> Dmem<span class="op">[</span>SIZE_DM<span class="dv">-1</span><span class="op">:</span><span class="dv">0</span><span class="op">];</span>      <span class="co">// Data memory array</span></span>
<span id="cb76-19"><a href="#cb76-19" aria-hidden="true" tabindex="-1"></a>  <span class="dt">integer</span> i<span class="op">;</span></span>
<span id="cb76-20"><a href="#cb76-20" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Initialize the data memory</span></span>
<span id="cb76-21"><a href="#cb76-21" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb76-22"><a href="#cb76-22" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Fill the data memory with zeros</span></span>
<span id="cb76-23"><a href="#cb76-23" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> SIZE_DM<span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb76-24"><a href="#cb76-24" aria-hidden="true" tabindex="-1"></a>      Dmem<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb76-25"><a href="#cb76-25" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb76-26"><a href="#cb76-26" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb76-27"><a href="#cb76-27" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Memory read operation</span></span>
<span id="cb76-28"><a href="#cb76-28" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_addr <span class="dt">or</span> i_MemRead <span class="dt">or</span> i_MemtoReg <span class="dt">or</span> i_ALUresult<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb76-29"><a href="#cb76-29" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_MemRead <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span>                  <span class="co">// If memory read is enabled</span></span>
<span id="cb76-30"><a href="#cb76-30" aria-hidden="true" tabindex="-1"></a>      <span class="kw">if</span> <span class="op">(</span>i_MemtoReg <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span>               <span class="co">// If MemtoReg is 1, select memory data as output</span></span>
<span id="cb76-31"><a href="#cb76-31" aria-hidden="true" tabindex="-1"></a>        o_rData <span class="op">=</span> Dmem<span class="op">[</span>i_addr<span class="op">];</span>                <span class="co">// Read data from the memory array</span></span>
<span id="cb76-32"><a href="#cb76-32" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb76-33"><a href="#cb76-33" aria-hidden="true" tabindex="-1"></a>        o_rData <span class="op">=</span> i_ALUresult<span class="op">;</span>                 <span class="co">// If MemtoReg is 0, select ALU result as output</span></span>
<span id="cb76-34"><a href="#cb76-34" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb76-35"><a href="#cb76-35" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb76-36"><a href="#cb76-36" aria-hidden="true" tabindex="-1"></a>      o_rData <span class="op">=</span> i_ALUresult<span class="op">;</span>                   <span class="co">// If memory read is not enabled, select ALU result as output</span></span>
<span id="cb76-37"><a href="#cb76-37" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb76-38"><a href="#cb76-38" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb76-39"><a href="#cb76-39" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Memory write operation</span></span>
<span id="cb76-40"><a href="#cb76-40" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_clk<span class="op">)</span> <span class="kw">begin</span>                <span class="co">// Triggered on the positive edge of the clock</span></span>
<span id="cb76-41"><a href="#cb76-41" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_MemWrite <span class="op">==</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span>                 <span class="co">// If memory write is enabled</span></span>
<span id="cb76-42"><a href="#cb76-42" aria-hidden="true" tabindex="-1"></a>      Dmem<span class="op">[</span>i_addr<span class="op">]</span> <span class="op">=</span> i_wData<span class="op">;</span>                  <span class="co">// Write data to the memory array</span></span>
<span id="cb76-43"><a href="#cb76-43" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb76-44"><a href="#cb76-44" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb76-45"><a href="#cb76-45" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Interactions with other components: - The <code>DataMemory</code>
module receives the address (<code>i_addr</code>), write data
(<code>i_wData</code>), and control signals (<code>i_MemWrite</code>,
<code>i_MemRead</code>, <code>i_MemtoReg</code>) from the
<code>ControlUnit</code> and <code>ALU</code> modules. - It provides the
read data (<code>o_rData</code>) to the <code>RegisterFile</code> module
for store instructions or to the <code>ALU</code> for load instructions.
- The <code>i_ALUresult</code> input is used as the memory address for
read and write operations. - The <code>i_MemWrite</code> control signal
determines whether a memory write operation should be performed. - The
<code>i_MemRead</code> control signal determines whether a memory read
operation should be performed. - The <code>i_MemtoReg</code> control
signal selects whether the memory data or the ALU result should be
output as the read data.</p>
<p>The <code>DataMemory</code> module is critical for providing data
storage and handling memory read and write operations.</p>
<p>It interacts with the control unit, ALU, and register file to
facilitate data movement and manipulation in the processor.</p>
<h3 id="instruction-memory-2">Instruction Memory</h3>
<p>Here the <code>ProgramCounter.v</code> file with detailed code
comments explaining its purpose, functionality, and interactions with
other components in the MIPS processor:</p>
<div class="sourceCode" id="cb77"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb77-1"><a href="#cb77-1" aria-hidden="true" tabindex="-1"></a><span class="co">// File: ProgramCounter.v</span></span>
<span id="cb77-2"><a href="#cb77-2" aria-hidden="true" tabindex="-1"></a><span class="co">// Description: This file contains the program counter module for the MIPS processor.</span></span>
<span id="cb77-3"><a href="#cb77-3" aria-hidden="true" tabindex="-1"></a><span class="co">// Purpose: The program counter keeps track of the current instruction address and updates it</span></span>
<span id="cb77-4"><a href="#cb77-4" aria-hidden="true" tabindex="-1"></a><span class="co">//          to the next instruction address on each clock cycle. It is responsible for providing</span></span>
<span id="cb77-5"><a href="#cb77-5" aria-hidden="true" tabindex="-1"></a><span class="co">//          the address of the instruction to be fetched from the instruction memory.</span></span>
<span id="cb77-6"><a href="#cb77-6" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb77-7"><a href="#cb77-7" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ProgramCounter <span class="op">(</span></span>
<span id="cb77-8"><a href="#cb77-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_Clk<span class="op">,</span>                <span class="co">// Input clock signal</span></span>
<span id="cb77-9"><a href="#cb77-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Next<span class="op">,</span>        <span class="co">// Input next instruction address</span></span>
<span id="cb77-10"><a href="#cb77-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_Out     <span class="co">// Output current instruction address</span></span>
<span id="cb77-11"><a href="#cb77-11" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb77-12"><a href="#cb77-12" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Initialize the program counter</span></span>
<span id="cb77-13"><a href="#cb77-13" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb77-14"><a href="#cb77-14" aria-hidden="true" tabindex="-1"></a>    o_Out <span class="op">=</span> <span class="op">-</span><span class="dv">4</span><span class="op">;</span>                 <span class="co">// Set the initial address to -4 (used for reset or initialization)</span></span>
<span id="cb77-15"><a href="#cb77-15" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb77-16"><a href="#cb77-16" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Update the program counter on the positive edge of the clock</span></span>
<span id="cb77-17"><a href="#cb77-17" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_Clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb77-18"><a href="#cb77-18" aria-hidden="true" tabindex="-1"></a>    o_Out <span class="op">=</span> i_Next<span class="op">;</span>             <span class="co">// Update the current address with the next address</span></span>
<span id="cb77-19"><a href="#cb77-19" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb77-20"><a href="#cb77-20" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Interactions with other components: - The <code>ProgramCounter</code>
module receives the next instruction address (<code>i_Next</code>) from
the <code>NextProgramCounter</code> module. - It provides the current
instruction address (<code>o_Out</code>) to the
<code>InstructionMemory</code> module to fetch the corresponding
instruction. - The <code>ProgramCounter</code> is updated on the
positive edge of the clock signal (<code>i_Clk</code>), which is
typically connected to the global clock signal of the processor.</p>
<p>The <code>ProgramCounter</code> module is a critical component in the
MIPS processor pipeline. It keeps track of the current instruction
address and updates it on each clock cycle to fetch the next
instruction. The program counter ensures the sequential execution of
instructions and enables the processor to navigate through the
program.</p>
<h3 id="program-counter-1">Program Counter</h3>
<p>Here the <code>ProgramCounter.v</code> file with detailed code
comments explaining its purpose, functionality, and interactions with
other components in the MIPS processor:</p>
<div class="sourceCode" id="cb78"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb78-1"><a href="#cb78-1" aria-hidden="true" tabindex="-1"></a><span class="co">// File: ProgramCounter.v</span></span>
<span id="cb78-2"><a href="#cb78-2" aria-hidden="true" tabindex="-1"></a><span class="co">// Description: This file contains the program counter module for the MIPS processor.</span></span>
<span id="cb78-3"><a href="#cb78-3" aria-hidden="true" tabindex="-1"></a><span class="co">// Purpose: The program counter keeps track of the current instruction address and updates it to the next address.</span></span>
<span id="cb78-4"><a href="#cb78-4" aria-hidden="true" tabindex="-1"></a><span class="co">//          It is responsible for providing the current instruction address to the instruction memory and updating</span></span>
<span id="cb78-5"><a href="#cb78-5" aria-hidden="true" tabindex="-1"></a><span class="co">//          the address based on the next address input.</span></span>
<span id="cb78-6"><a href="#cb78-6" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb78-7"><a href="#cb78-7" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ProgramCounter <span class="op">(</span></span>
<span id="cb78-8"><a href="#cb78-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> i_Clk<span class="op">,</span>                   <span class="co">// Input clock signal</span></span>
<span id="cb78-9"><a href="#cb78-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Next<span class="op">,</span>           <span class="co">// Input next instruction address</span></span>
<span id="cb78-10"><a href="#cb78-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_Out        <span class="co">// Output current instruction address</span></span>
<span id="cb78-11"><a href="#cb78-11" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb78-12"><a href="#cb78-12" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Initialize the program counter</span></span>
<span id="cb78-13"><a href="#cb78-13" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb78-14"><a href="#cb78-14" aria-hidden="true" tabindex="-1"></a>    o_Out <span class="op">=</span> <span class="op">-</span><span class="dv">4</span><span class="op">;</span>                    <span class="co">// Set the initial instruction address to -4</span></span>
<span id="cb78-15"><a href="#cb78-15" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb78-16"><a href="#cb78-16" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Update the program counter on the positive edge of the clock</span></span>
<span id="cb78-17"><a href="#cb78-17" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span><span class="kw">posedge</span> i_Clk<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb78-18"><a href="#cb78-18" aria-hidden="true" tabindex="-1"></a>    o_Out <span class="op">=</span> i_Next<span class="op">;</span>                <span class="co">// Update the current instruction address with the next address</span></span>
<span id="cb78-19"><a href="#cb78-19" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb78-20"><a href="#cb78-20" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Purpose and Functionality: - The <code>ProgramCounter</code> module
keeps track of the current instruction address in the MIPS processor. -
It is responsible for providing the current instruction address to the
instruction memory (<code>InstructionMemory</code>) for fetching the
corresponding instruction. - The program counter is updated on the
positive edge of the clock signal (<code>i_Clk</code>). - The next
instruction address (<code>i_Next</code>) is provided as an input to the
module, which is used to update the current instruction address
(<code>o_Out</code>) on each clock cycle. - The initial value of the
program counter is set to -4, which represents the initial state before
the first instruction is fetched.</p>
<p>Interactions with other components: - The <code>ProgramCounter</code>
module receives the next instruction address (<code>i_Next</code>) from
the <code>NextProgramCounter</code> module, which calculates the next
address based on the current instruction and control signals. - It
provides the current instruction address (<code>o_Out</code>) to the
<code>InstructionMemory</code> module to fetch the corresponding
instruction. - The <code>ProgramCounter</code> is updated on the
positive edge of the clock signal (<code>i_Clk</code>), which is
typically connected to the main processor clock.</p>
<p>The <code>ProgramCounter</code> module is to manage the flow of
execution by keeping track of the current instruction address.</p>
<p>It ensures that instructions are fetched and executed in the correct
order by updating the address on each clock cycle based on the next
address input provided by the <code>NextProgramCounter</code>
module.</p>
<h3 id="alu-2">ALU</h3>
<p>Here the <code>ALU.v</code> file with detailed code comments
explaining its purpose, functionality, and interactions with other
components in the MIPS processor:</p>
<div class="sourceCode" id="cb79"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb79-1"><a href="#cb79-1" aria-hidden="true" tabindex="-1"></a><span class="co">// File: ALU.v</span></span>
<span id="cb79-2"><a href="#cb79-2" aria-hidden="true" tabindex="-1"></a><span class="co">// Description: This file contains the Arithmetic Logic Unit (ALU) module for the MIPS processor.</span></span>
<span id="cb79-3"><a href="#cb79-3" aria-hidden="true" tabindex="-1"></a><span class="co">// Purpose: The ALU performs arithmetic and logic operations based on the ALU control signals.</span></span>
<span id="cb79-4"><a href="#cb79-4" aria-hidden="true" tabindex="-1"></a><span class="co">//          It takes two input operands (i_data1 and i_read2/immediate value) and performs the specified operation.</span></span>
<span id="cb79-5"><a href="#cb79-5" aria-hidden="true" tabindex="-1"></a><span class="co">//          The ALU also generates a zero flag (o_Zero) to indicate if the result is zero.</span></span>
<span id="cb79-6"><a href="#cb79-6" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb79-7"><a href="#cb79-7" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb79-8"><a href="#cb79-8" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb79-9"><a href="#cb79-9" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ALU <span class="op">(</span></span>
<span id="cb79-10"><a href="#cb79-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_data1<span class="op">,</span>        <span class="co">// Input operand 1 (from RegisterFile)</span></span>
<span id="cb79-11"><a href="#cb79-11" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_read2<span class="op">,</span>        <span class="co">// Input operand 2 (from RegisterFile or immediate value)</span></span>
<span id="cb79-12"><a href="#cb79-12" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_Instruction<span class="op">,</span>  <span class="co">// Input instruction (used for sign-extension of immediate value)</span></span>
<span id="cb79-13"><a href="#cb79-13" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>             i_ALUSrc<span class="op">,</span>       <span class="co">// Control signal to select between i_read2 or immediate value</span></span>
<span id="cb79-14"><a href="#cb79-14" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span>      <span class="op">[</span> <span class="dv">3</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_ALUcontrol<span class="op">,</span>   <span class="co">// Control signal to specify the ALU operation</span></span>
<span id="cb79-15"><a href="#cb79-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span>        o_Zero<span class="op">,</span>         <span class="co">// Output zero flag (1 if the ALU result is zero, 0 otherwise)</span></span>
<span id="cb79-16"><a href="#cb79-16" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_ALUresult     <span class="co">// Output ALU result</span></span>
<span id="cb79-17"><a href="#cb79-17" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb79-18"><a href="#cb79-18" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> data2<span class="op">;</span></span>
<span id="cb79-19"><a href="#cb79-19" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Determine the second operand based on the ALUSrc control signal</span></span>
<span id="cb79-20"><a href="#cb79-20" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_ALUSrc<span class="op">,</span> i_read2<span class="op">,</span> i_Instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb79-21"><a href="#cb79-21" aria-hidden="true" tabindex="-1"></a>    <span class="kw">if</span> <span class="op">(</span>i_ALUSrc <span class="op">==</span> <span class="dv">0</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb79-22"><a href="#cb79-22" aria-hidden="true" tabindex="-1"></a>      data2 <span class="op">=</span> i_read2<span class="op">;</span>                 <span class="co">// Use i_read2 as the second operand</span></span>
<span id="cb79-23"><a href="#cb79-23" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb79-24"><a href="#cb79-24" aria-hidden="true" tabindex="-1"></a>      <span class="co">// Sign-extend the immediate value</span></span>
<span id="cb79-25"><a href="#cb79-25" aria-hidden="true" tabindex="-1"></a>      <span class="kw">if</span> <span class="op">(</span>i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">]</span> <span class="op">==</span> <span class="bn">1&#39;b0</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb79-26"><a href="#cb79-26" aria-hidden="true" tabindex="-1"></a>        data2 <span class="op">=</span> <span class="op">{</span><span class="bn">16&#39;b0</span><span class="op">,</span> i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]};</span>  <span class="co">// Zero-extend if the immediate value is positive</span></span>
<span id="cb79-27"><a href="#cb79-27" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span> <span class="kw">else</span> <span class="kw">begin</span></span>
<span id="cb79-28"><a href="#cb79-28" aria-hidden="true" tabindex="-1"></a>        data2 <span class="op">=</span> <span class="op">{{</span><span class="dv">16</span><span class="op">{</span><span class="bn">1&#39;b1</span><span class="op">}},</span> i_Instruction<span class="op">[</span><span class="dv">15</span><span class="op">:</span><span class="dv">0</span><span class="op">]};</span>  <span class="co">// Sign-extend if the immediate value is negative</span></span>
<span id="cb79-29"><a href="#cb79-29" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb79-30"><a href="#cb79-30" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb79-31"><a href="#cb79-31" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb79-32"><a href="#cb79-32" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Perform the ALU operation based on the ALUcontrol signal</span></span>
<span id="cb79-33"><a href="#cb79-33" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> <span class="op">@(</span>i_data1<span class="op">,</span> data2<span class="op">,</span> i_ALUcontrol<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb79-34"><a href="#cb79-34" aria-hidden="true" tabindex="-1"></a>    <span class="kw">case</span> <span class="op">(</span>i_ALUcontrol<span class="op">)</span></span>
<span id="cb79-35"><a href="#cb79-35" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0000</span><span class="op">:</span>  <span class="co">// AND</span></span>
<span id="cb79-36"><a href="#cb79-36" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> i_data1 <span class="op">&amp;</span> data2<span class="op">;</span></span>
<span id="cb79-37"><a href="#cb79-37" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0001</span><span class="op">:</span>  <span class="co">// OR</span></span>
<span id="cb79-38"><a href="#cb79-38" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> i_data1 <span class="op">|</span> data2<span class="op">;</span></span>
<span id="cb79-39"><a href="#cb79-39" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0010</span><span class="op">:</span>  <span class="co">// ADD</span></span>
<span id="cb79-40"><a href="#cb79-40" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> i_data1 <span class="op">+</span> data2<span class="op">;</span></span>
<span id="cb79-41"><a href="#cb79-41" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0110</span><span class="op">:</span>  <span class="co">// SUB</span></span>
<span id="cb79-42"><a href="#cb79-42" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> i_data1 <span class="op">-</span> data2<span class="op">;</span></span>
<span id="cb79-43"><a href="#cb79-43" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b0111</span><span class="op">:</span>  <span class="co">// SLT (Set Less Than)</span></span>
<span id="cb79-44"><a href="#cb79-44" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> <span class="op">(</span>i_data1 <span class="op">&lt;</span> data2<span class="op">)</span> <span class="op">?</span> <span class="dv">1</span> <span class="op">:</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb79-45"><a href="#cb79-45" aria-hidden="true" tabindex="-1"></a>      <span class="bn">4&#39;b1100</span><span class="op">:</span>  <span class="co">// NOR</span></span>
<span id="cb79-46"><a href="#cb79-46" aria-hidden="true" tabindex="-1"></a>        o_ALUresult <span class="op">=</span> <span class="op">~(</span>i_data1 <span class="op">|</span> data2<span class="op">);</span></span>
<span id="cb79-47"><a href="#cb79-47" aria-hidden="true" tabindex="-1"></a>      <span class="kw">default</span><span class="op">:</span> <span class="op">;</span></span>
<span id="cb79-48"><a href="#cb79-48" aria-hidden="true" tabindex="-1"></a>    <span class="kw">endcase</span></span>
<span id="cb79-49"><a href="#cb79-49" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Set the zero flag if the ALU result is zero</span></span>
<span id="cb79-50"><a href="#cb79-50" aria-hidden="true" tabindex="-1"></a>    o_Zero <span class="op">=</span> <span class="op">(</span>o_ALUresult <span class="op">==</span> <span class="dv">0</span><span class="op">)</span> <span class="op">?</span> <span class="dv">1</span> <span class="op">:</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb79-51"><a href="#cb79-51" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb79-52"><a href="#cb79-52" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb79-53"><a href="#cb79-53" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Interactions with other components: - The <code>ALU</code> module
receives input operands (<code>i_data1</code> and <code>i_read2</code>)
from the <code>RegisterFile</code> module. - The <code>i_ALUSrc</code>
control signal from the <code>ControlUnit</code> determines whether the
second operand is <code>i_read2</code> or an immediate value from the
instruction (<code>i_Instruction</code>). - The
<code>i_ALUcontrol</code> signal from the <code>ALUControl</code> module
specifies the ALU operation to be performed. - The <code>ALU</code>
module outputs the result (<code>o_ALUresult</code>) to the
<code>DataMemory</code> and <code>RegisterFile</code> modules for memory
access and register writeback. - The zero flag (<code>o_Zero</code>) is
used by</p>
<h3 id="control-unit-1">Control Unit</h3>
<p>Certainly! Here’s a detailed explanation of the Verilog module
provided for a Control Unit in a single cycle MIPS processor. Each line
of the module is annotated to explain its function and relevance.</p>
<div class="sourceCode" id="cb80"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb80-1"><a href="#cb80-1" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb80-2"><a href="#cb80-2" aria-hidden="true" tabindex="-1"></a><span class="co">// Defines the time unit as 1 nanosecond and the simulation time precision as 1 picosecond.</span></span>
<span id="cb80-3"><a href="#cb80-3" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> ControlUnit <span class="op">(</span></span>
<span id="cb80-4"><a href="#cb80-4" aria-hidden="true" tabindex="-1"></a>    <span class="dt">input</span> <span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> i_instruction<span class="op">,</span>     <span class="co">// 32-bit input for the instruction.</span></span>
<span id="cb80-5"><a href="#cb80-5" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_RegDst<span class="op">,</span>            <span class="co">// Determines if rd (1) or rt (0) should be the destination register.</span></span>
<span id="cb80-6"><a href="#cb80-6" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_Jump<span class="op">,</span>              <span class="co">// Control signal for jumping to an instruction address.</span></span>
<span id="cb80-7"><a href="#cb80-7" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_Branch<span class="op">,</span>            <span class="co">// Control signal for branching (beq).</span></span>
<span id="cb80-8"><a href="#cb80-8" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_Bne<span class="op">,</span>               <span class="co">// Control signal for branching not equal (bne).</span></span>
<span id="cb80-9"><a href="#cb80-9" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_MemRead<span class="op">,</span>           <span class="co">// Enables reading from memory (used by lw).</span></span>
<span id="cb80-10"><a href="#cb80-10" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_MemtoReg<span class="op">,</span>          <span class="co">// Determines if the value should come from memory (1) or ALU (0).</span></span>
<span id="cb80-11"><a href="#cb80-11" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">1</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_ALUOp<span class="op">,</span>       <span class="co">// Control signal for ALU operation type.</span></span>
<span id="cb80-12"><a href="#cb80-12" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_MemWrite<span class="op">,</span>          <span class="co">// Enables writing to memory (used by sw).</span></span>
<span id="cb80-13"><a href="#cb80-13" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_ALUSrc<span class="op">,</span>            <span class="co">// Determines if the second ALU operand is an immediate (1) or register (0).</span></span>
<span id="cb80-14"><a href="#cb80-14" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> o_RegWrite<span class="op">,</span>          <span class="co">// Enables writing to the register file.</span></span>
<span id="cb80-15"><a href="#cb80-15" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_first<span class="op">,</span>   <span class="co">// Segment display outputs to visually represent instruction types or states.</span></span>
<span id="cb80-16"><a href="#cb80-16" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_second<span class="op">,</span>  <span class="co">// Each segment holds a 7-segment representation.</span></span>
<span id="cb80-17"><a href="#cb80-17" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_third<span class="op">,</span></span>
<span id="cb80-18"><a href="#cb80-18" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_fourth<span class="op">,</span></span>
<span id="cb80-19"><a href="#cb80-19" aria-hidden="true" tabindex="-1"></a>    <span class="dt">output</span> <span class="dt">reg</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> o_seg_fifth</span>
<span id="cb80-20"><a href="#cb80-20" aria-hidden="true" tabindex="-1"></a><span class="op">);</span></span>
<span id="cb80-21"><a href="#cb80-21" aria-hidden="true" tabindex="-1"></a><span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb80-22"><a href="#cb80-22" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize all control signals and display outputs to their default (usually disabled) states.</span></span>
<span id="cb80-23"><a href="#cb80-23" aria-hidden="true" tabindex="-1"></a>    o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-24"><a href="#cb80-24" aria-hidden="true" tabindex="-1"></a>    o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-25"><a href="#cb80-25" aria-hidden="true" tabindex="-1"></a>    o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-26"><a href="#cb80-26" aria-hidden="true" tabindex="-1"></a>    o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-27"><a href="#cb80-27" aria-hidden="true" tabindex="-1"></a>    o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-28"><a href="#cb80-28" aria-hidden="true" tabindex="-1"></a>    o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span> <span class="co">// Default ALU operation, no operation specified.</span></span>
<span id="cb80-29"><a href="#cb80-29" aria-hidden="true" tabindex="-1"></a>    o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-30"><a href="#cb80-30" aria-hidden="true" tabindex="-1"></a>    o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-31"><a href="#cb80-31" aria-hidden="true" tabindex="-1"></a>    o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-32"><a href="#cb80-32" aria-hidden="true" tabindex="-1"></a>    o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// All segments off (blank).</span></span>
<span id="cb80-33"><a href="#cb80-33" aria-hidden="true" tabindex="-1"></a>    o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span></span>
<span id="cb80-34"><a href="#cb80-34" aria-hidden="true" tabindex="-1"></a>    o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span></span>
<span id="cb80-35"><a href="#cb80-35" aria-hidden="true" tabindex="-1"></a>    o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span></span>
<span id="cb80-36"><a href="#cb80-36" aria-hidden="true" tabindex="-1"></a>    o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span></span>
<span id="cb80-37"><a href="#cb80-37" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span>
<span id="cb80-38"><a href="#cb80-38" aria-hidden="true" tabindex="-1"></a><span class="kw">always</span> <span class="op">@(</span>i_instruction<span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb80-39"><a href="#cb80-39" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Control logic triggered by any change in the instruction input.</span></span>
<span id="cb80-40"><a href="#cb80-40" aria-hidden="true" tabindex="-1"></a>    <span class="kw">case</span> <span class="op">(</span>i_instruction<span class="op">[</span><span class="dv">31</span><span class="op">:</span><span class="dv">26</span><span class="op">])</span> <span class="co">// Decode the opcode part of the instruction.</span></span>
<span id="cb80-41"><a href="#cb80-41" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000000</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// ARITHMETIC (R-type instructions)</span></span>
<span id="cb80-42"><a href="#cb80-42" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-43"><a href="#cb80-43" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-44"><a href="#cb80-44" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-45"><a href="#cb80-45" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-46"><a href="#cb80-46" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-47"><a href="#cb80-47" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-48"><a href="#cb80-48" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-49"><a href="#cb80-49" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-50"><a href="#cb80-50" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b10</span><span class="op">;</span> <span class="co">// Specific ALU operation for arithmetic.</span></span>
<span id="cb80-51"><a href="#cb80-51" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-52"><a href="#cb80-52" aria-hidden="true" tabindex="-1"></a>        <span class="co">// Display setup for ARITHMETIC.</span></span>
<span id="cb80-53"><a href="#cb80-53" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span>  <span class="bn">7&#39;b0001000</span><span class="op">;</span>  <span class="co">// A</span></span>
<span id="cb80-54"><a href="#cb80-54" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111010</span><span class="op">;</span>  <span class="co">// R</span></span>
<span id="cb80-55"><a href="#cb80-55" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span>  <span class="bn">7&#39;b1111001</span><span class="op">;</span>  <span class="co">// I</span></span>
<span id="cb80-56"><a href="#cb80-56" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b0001111</span><span class="op">;</span>  <span class="co">// T</span></span>
<span id="cb80-57"><a href="#cb80-57" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span>  <span class="bn">7&#39;b0001001</span><span class="op">;</span>  <span class="co">// H</span></span>
<span id="cb80-58"><a href="#cb80-58" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-59"><a href="#cb80-59" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b001000</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// addi</span></span>
<span id="cb80-60"><a href="#cb80-60" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-61"><a href="#cb80-61" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-62"><a href="#cb80-62" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-63"><a href="#cb80-63" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-64"><a href="#cb80-64" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-65"><a href="#cb80-65" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-66"><a href="#cb80-66" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-67"><a href="#cb80-67" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-68"><a href="#cb80-68" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb80-69"><a href="#cb80-69" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-70"><a href="#cb80-70" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b0001000</span><span class="op">;</span>  <span class="co">// A</span></span>
<span id="cb80-71"><a href="#cb80-71" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1000010</span><span class="op">;</span> <span class="co">// d</span></span>
<span id="cb80-72"><a href="#cb80-72" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1000010</span><span class="op">;</span>  <span class="co">// d</span></span>
<span id="cb80-73"><a href="#cb80-73" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1001111</span><span class="op">;</span> <span class="co">// i</span></span>
<span id="cb80-74"><a href="#cb80-74" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-75"><a href="#cb80-75" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-76"><a href="#cb80-76" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b001100</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// andi</span></span>
<span id="cb80-77"><a href="#cb80-77" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-78"><a href="#cb80-78" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-79"><a href="#cb80-79" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-80"><a href="#cb80-80" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-81"><a href="#cb80-81" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-82"><a href="#cb80-82" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-83"><a href="#cb80-83" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-84"><a href="#cb80-84" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-85"><a href="#cb80-85" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b11</span><span class="op">;</span></span>
<span id="cb80-86"><a href="#cb80-86" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-87"><a href="#cb80-87" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b0001000</span><span class="op">;</span>  <span class="co">// A</span></span>
<span id="cb80-88"><a href="#cb80-88" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b0101011</span><span class="op">;</span>  <span class="co">// n</span></span>
<span id="cb80-89"><a href="#cb80-89" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1000010</span><span class="op">;</span>  <span class="co">// d</span></span>
<span id="cb80-90"><a href="#cb80-90" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1001111</span><span class="op">;</span>  <span class="co">// i</span></span>
<span id="cb80-91"><a href="#cb80-91" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-92"><a href="#cb80-92" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-93"><a href="#cb80-93" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b100011</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// lw</span></span>
<span id="cb80-94"><a href="#cb80-94" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-95"><a href="#cb80-95" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-96"><a href="#cb80-96" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-97"><a href="#cb80-97" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-98"><a href="#cb80-98" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-99"><a href="#cb80-99" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-100"><a href="#cb80-100" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-101"><a href="#cb80-101" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-102"><a href="#cb80-102" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb80-103"><a href="#cb80-103" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-104"><a href="#cb80-104" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1000111</span><span class="op">;</span>  <span class="co">// L</span></span>
<span id="cb80-105"><a href="#cb80-105" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1001001</span><span class="op">;</span>  <span class="co">// w</span></span>
<span id="cb80-106"><a href="#cb80-106" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-107"><a href="#cb80-107" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-108"><a href="#cb80-108" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-109"><a href="#cb80-109" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-110"><a href="#cb80-110" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b101011</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// sw</span></span>
<span id="cb80-111"><a href="#cb80-111" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-112"><a href="#cb80-112" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-113"><a href="#cb80-113" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-114"><a href="#cb80-114" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-115"><a href="#cb80-115" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-116"><a href="#cb80-116" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-117"><a href="#cb80-117" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-118"><a href="#cb80-118" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-119"><a href="#cb80-119" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb80-120"><a href="#cb80-120" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-121"><a href="#cb80-121" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b0010010</span><span class="op">;</span>  <span class="co">// S</span></span>
<span id="cb80-122"><a href="#cb80-122" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1001001</span><span class="op">;</span>  <span class="co">// w</span></span>
<span id="cb80-123"><a href="#cb80-123" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-124"><a href="#cb80-124" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-125"><a href="#cb80-125" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-126"><a href="#cb80-126" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-127"><a href="#cb80-127" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000100</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// beq</span></span>
<span id="cb80-128"><a href="#cb80-128" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-129"><a href="#cb80-129" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-130"><a href="#cb80-130" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-131"><a href="#cb80-131" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-132"><a href="#cb80-132" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-133"><a href="#cb80-133" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-134"><a href="#cb80-134" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-135"><a href="#cb80-135" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-136"><a href="#cb80-136" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b01</span><span class="op">;</span></span>
<span id="cb80-137"><a href="#cb80-137" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-138"><a href="#cb80-138" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1100000</span><span class="op">;</span>  <span class="co">// b</span></span>
<span id="cb80-139"><a href="#cb80-139" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b0110000</span><span class="op">;</span>  <span class="co">// e</span></span>
<span id="cb80-140"><a href="#cb80-140" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b0001100</span><span class="op">;</span>  <span class="co">// q</span></span>
<span id="cb80-141"><a href="#cb80-141" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-142"><a href="#cb80-142" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-143"><a href="#cb80-143" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-144"><a href="#cb80-144" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000101</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// bne</span></span>
<span id="cb80-145"><a href="#cb80-145" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-146"><a href="#cb80-146" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-147"><a href="#cb80-147" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-148"><a href="#cb80-148" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-149"><a href="#cb80-149" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-150"><a href="#cb80-150" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-151"><a href="#cb80-151" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-152"><a href="#cb80-152" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-153"><a href="#cb80-153" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b01</span><span class="op">;</span></span>
<span id="cb80-154"><a href="#cb80-154" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-155"><a href="#cb80-155" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1100000</span><span class="op">;</span>  <span class="co">// b</span></span>
<span id="cb80-156"><a href="#cb80-156" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b0101011</span><span class="op">;</span> <span class="co">// n</span></span>
<span id="cb80-157"><a href="#cb80-157" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b0110000</span><span class="op">;</span>  <span class="co">// e</span></span>
<span id="cb80-158"><a href="#cb80-158" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span> <span class="co">// Blank</span></span>
<span id="cb80-159"><a href="#cb80-159" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-160"><a href="#cb80-160" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-161"><a href="#cb80-161" aria-hidden="true" tabindex="-1"></a>      <span class="bn">6&#39;b000010</span><span class="op">:</span> <span class="kw">begin</span>  <span class="co">// j</span></span>
<span id="cb80-162"><a href="#cb80-162" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-163"><a href="#cb80-163" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-164"><a href="#cb80-164" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>  <span class="co">// X</span></span>
<span id="cb80-165"><a href="#cb80-165" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-166"><a href="#cb80-166" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-167"><a href="#cb80-167" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-168"><a href="#cb80-168" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-169"><a href="#cb80-169" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-170"><a href="#cb80-170" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b01</span><span class="op">;</span></span>
<span id="cb80-171"><a href="#cb80-171" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">1</span><span class="op">;</span></span>
<span id="cb80-172"><a href="#cb80-172" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1100001</span><span class="op">;</span>  <span class="co">// J</span></span>
<span id="cb80-173"><a href="#cb80-173" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span> <span class="co">// Blank</span></span>
<span id="cb80-174"><a href="#cb80-174" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-175"><a href="#cb80-175" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span> <span class="co">// Blank</span></span>
<span id="cb80-176"><a href="#cb80-176" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-177"><a href="#cb80-177" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-178"><a href="#cb80-178" aria-hidden="true" tabindex="-1"></a>      <span class="kw">default</span><span class="op">:</span> <span class="kw">begin</span></span>
<span id="cb80-179"><a href="#cb80-179" aria-hidden="true" tabindex="-1"></a>        <span class="co">// Default case sets all outputs to zero or disables them, providing a safe default state.</span></span>
<span id="cb80-180"><a href="#cb80-180" aria-hidden="true" tabindex="-1"></a>        o_RegDst <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-181"><a href="#cb80-181" aria-hidden="true" tabindex="-1"></a>        o_ALUSrc <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-182"><a href="#cb80-182" aria-hidden="true" tabindex="-1"></a>        o_MemtoReg <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-183"><a href="#cb80-183" aria-hidden="true" tabindex="-1"></a>        o_RegWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-184"><a href="#cb80-184" aria-hidden="true" tabindex="-1"></a>        o_MemRead <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-185"><a href="#cb80-185" aria-hidden="true" tabindex="-1"></a>        o_MemWrite <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-186"><a href="#cb80-186" aria-hidden="true" tabindex="-1"></a>        o_Branch <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-187"><a href="#cb80-187" aria-hidden="true" tabindex="-1"></a>        o_Bne <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-188"><a href="#cb80-188" aria-hidden="true" tabindex="-1"></a>        o_ALUOp <span class="op">=</span> <span class="bn">2&#39;b00</span><span class="op">;</span></span>
<span id="cb80-189"><a href="#cb80-189" aria-hidden="true" tabindex="-1"></a>        o_Jump <span class="op">=</span> <span class="dv">0</span><span class="op">;</span></span>
<span id="cb80-190"><a href="#cb80-190" aria-hidden="true" tabindex="-1"></a>        <span class="co">// Display all segments off for undefined instructions.</span></span>
<span id="cb80-191"><a href="#cb80-191" aria-hidden="true" tabindex="-1"></a>        o_seg_first <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-192"><a href="#cb80-192" aria-hidden="true" tabindex="-1"></a>        o_seg_second <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span> <span class="co">// Blank</span></span>
<span id="cb80-193"><a href="#cb80-193" aria-hidden="true" tabindex="-1"></a>        o_seg_third <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-194"><a href="#cb80-194" aria-hidden="true" tabindex="-1"></a>        o_seg_fourth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span> <span class="co">// Blank</span></span>
<span id="cb80-195"><a href="#cb80-195" aria-hidden="true" tabindex="-1"></a>        o_seg_fifth <span class="op">=</span> <span class="bn">7&#39;b1111111</span><span class="op">;</span>  <span class="co">// Blank</span></span>
<span id="cb80-196"><a href="#cb80-196" aria-hidden="true" tabindex="-1"></a>      <span class="kw">end</span></span>
<span id="cb80-197"><a href="#cb80-197" aria-hidden="true" tabindex="-1"></a>    <span class="kw">endcase</span></span>
<span id="cb80-198"><a href="#cb80-198" aria-hidden="true" tabindex="-1"></a><span class="kw">end</span></span>
<span id="cb80-199"><a href="#cb80-199" aria-hidden="true" tabindex="-1"></a></span>
<span id="cb80-200"><a href="#cb80-200" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>This code serves as the control logic for a single cycle MIPS
processor, managing the routing and operations of data based on the
instruction being executed. It adjusts the path and operation of the
data in various parts of the processor according to the opcode of the
instruction, with added visual output for debugging or educational
purposes through a 7-segment display configuration.</p>
<h3 id="testbench">Testbench</h3>
<p>Here is the <code>mips_tb.v</code> file with detailed code comments
explaining its purpose, functionality, and interactions with other
components in the MIPS processor:</p>
<div class="sourceCode" id="cb81"><pre
class="sourceCode verilog"><code class="sourceCode verilog"><span id="cb81-1"><a href="#cb81-1" aria-hidden="true" tabindex="-1"></a><span class="co">// File: mips_tb.v</span></span>
<span id="cb81-2"><a href="#cb81-2" aria-hidden="true" tabindex="-1"></a><span class="co">// Description: This file contains the testbench for the MIPS processor.</span></span>
<span id="cb81-3"><a href="#cb81-3" aria-hidden="true" tabindex="-1"></a><span class="co">// Purpose: The testbench is used to simulate and verify the functionality of the MIPS processor.</span></span>
<span id="cb81-4"><a href="#cb81-4" aria-hidden="true" tabindex="-1"></a><span class="co">//          It instantiates the MIPS processor module, provides clock and reset signals, and</span></span>
<span id="cb81-5"><a href="#cb81-5" aria-hidden="true" tabindex="-1"></a><span class="co">//          initializes the data memory and register file. It also displays the output on 7-segment displays.</span></span>
<span id="cb81-6"><a href="#cb81-6" aria-hidden="true" tabindex="-1"></a><span class="ot">`timescale 1ns / 1ps</span></span>
<span id="cb81-7"><a href="#cb81-7" aria-hidden="true" tabindex="-1"></a><span class="ot">`define CYCLE_TIME 20</span></span>
<span id="cb81-8"><a href="#cb81-8" aria-hidden="true" tabindex="-1"></a><span class="kw">module</span> mips_tb<span class="op">;</span></span>
<span id="cb81-9"><a href="#cb81-9" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> clk<span class="op">;</span>                 <span class="co">// Clock signal</span></span>
<span id="cb81-10"><a href="#cb81-10" aria-hidden="true" tabindex="-1"></a>  <span class="dt">reg</span> rst<span class="op">;</span>                 <span class="co">// Reset signal</span></span>
<span id="cb81-11"><a href="#cb81-11" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Segments for the 7-segment displays</span></span>
<span id="cb81-12"><a href="#cb81-12" aria-hidden="true" tabindex="-1"></a>  <span class="dt">wire</span> <span class="op">[</span><span class="dv">6</span><span class="op">:</span><span class="dv">0</span><span class="op">]</span> seg_first<span class="op">,</span> seg_second<span class="op">,</span> seg_third<span class="op">,</span> seg_fourth<span class="op">,</span> seg_fifth<span class="op">;</span></span>
<span id="cb81-13"><a href="#cb81-13" aria-hidden="true" tabindex="-1"></a>  <span class="dt">integer</span> i<span class="op">;</span>               <span class="co">// Loop variable</span></span>
<span id="cb81-14"><a href="#cb81-14" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Generate clock signal</span></span>
<span id="cb81-15"><a href="#cb81-15" aria-hidden="true" tabindex="-1"></a>  <span class="kw">always</span> #<span class="op">(</span><span class="ot">`CYCLE_TIME</span> <span class="op">/</span> <span class="dv">2</span><span class="op">)</span> clk <span class="op">=</span> <span class="op">~</span>clk<span class="op">;</span></span>
<span id="cb81-16"><a href="#cb81-16" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Instantiate the MIPS processor module</span></span>
<span id="cb81-17"><a href="#cb81-17" aria-hidden="true" tabindex="-1"></a>  mips uut <span class="op">(</span></span>
<span id="cb81-18"><a href="#cb81-18" aria-hidden="true" tabindex="-1"></a>      .i_Clk<span class="op">(</span>clk<span class="op">),</span></span>
<span id="cb81-19"><a href="#cb81-19" aria-hidden="true" tabindex="-1"></a>      .i_Rst<span class="op">(</span>rst<span class="op">),</span></span>
<span id="cb81-20"><a href="#cb81-20" aria-hidden="true" tabindex="-1"></a>      .o_Seg_first<span class="op">(</span>seg_first<span class="op">),</span></span>
<span id="cb81-21"><a href="#cb81-21" aria-hidden="true" tabindex="-1"></a>      .o_Seg_second<span class="op">(</span>seg_second<span class="op">),</span></span>
<span id="cb81-22"><a href="#cb81-22" aria-hidden="true" tabindex="-1"></a>      .o_Seg_third<span class="op">(</span>seg_third<span class="op">),</span></span>
<span id="cb81-23"><a href="#cb81-23" aria-hidden="true" tabindex="-1"></a>      .o_Seg_fourth<span class="op">(</span>seg_fourth<span class="op">),</span></span>
<span id="cb81-24"><a href="#cb81-24" aria-hidden="true" tabindex="-1"></a>      .o_Seg_fifth<span class="op">(</span>seg_fifth<span class="op">)</span></span>
<span id="cb81-25"><a href="#cb81-25" aria-hidden="true" tabindex="-1"></a>  <span class="op">);</span></span>
<span id="cb81-26"><a href="#cb81-26" aria-hidden="true" tabindex="-1"></a>  <span class="co">// Initialize data memory and register file</span></span>
<span id="cb81-27"><a href="#cb81-27" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb81-28"><a href="#cb81-28" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize data memory</span></span>
<span id="cb81-29"><a href="#cb81-29" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb81-30"><a href="#cb81-30" aria-hidden="true" tabindex="-1"></a>      uut.inst_DataMemory.Dmem<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb81-31"><a href="#cb81-31" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb81-32"><a href="#cb81-32" aria-hidden="true" tabindex="-1"></a>    <span class="co">// Initialize register file</span></span>
<span id="cb81-33"><a href="#cb81-33" aria-hidden="true" tabindex="-1"></a>    <span class="kw">for</span> <span class="op">(</span>i <span class="op">=</span> <span class="dv">0</span><span class="op">;</span> i <span class="op">&lt;</span> <span class="dv">32</span><span class="op">;</span> i <span class="op">=</span> i <span class="op">+</span> <span class="dv">1</span><span class="op">)</span> <span class="kw">begin</span></span>
<span id="cb81-34"><a href="#cb81-34" aria-hidden="true" tabindex="-1"></a>      uut.inst_RegisterFile.RegData<span class="op">[</span>i<span class="op">]</span> <span class="op">=</span> <span class="bn">32&#39;b0</span><span class="op">;</span></span>
<span id="cb81-35"><a href="#cb81-35" aria-hidden="true" tabindex="-1"></a>    <span class="kw">end</span></span>
<span id="cb81-36"><a href="#cb81-36" aria-hidden="true" tabindex="-1"></a>    clk <span class="op">=</span> <span class="dv">0</span><span class="op">;</span>                <span class="co">// Initialize clock signal</span></span>
<span id="cb81-37"><a href="#cb81-37" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb81-38"><a href="#cb81-38" aria-hidden="true" tabindex="-1"></a>  <span class="kw">initial</span> <span class="kw">begin</span></span>
<span id="cb81-39"><a href="#cb81-39" aria-hidden="true" tabindex="-1"></a>    <span class="bn">#1800</span> <span class="dt">$finish</span><span class="op">;</span></span>
<span id="cb81-40"><a href="#cb81-40" aria-hidden="true" tabindex="-1"></a>  <span class="kw">end</span></span>
<span id="cb81-41"><a href="#cb81-41" aria-hidden="true" tabindex="-1"></a><span class="kw">endmodule</span></span></code></pre></div>
<p>Interactions with other components:</p>
<ul>
<li>The <code>mips_tb</code> module instantiates the <code>mips</code>
module, which represents the MIPS processor.</li>
<li>It provides the clock signal (<code>clk</code>) to the
<code>mips</code> module for synchronization.</li>
<li>The reset signal (<code>rst</code>) is not used in this testbench
but can be used to reset the processor if needed.</li>
<li>The testbench initializes the data memory
(<code>inst_DataMemory.Dmem</code>) and register file
(<code>inst_RegisterFile.RegData</code>) of the <code>mips</code> module
to zero.</li>
<li>The 7-segment display outputs (<code>seg_first</code>,
<code>seg_second</code>, <code>seg_third</code>,
<code>seg_fourth</code>, <code>seg_fifth</code>) from the
<code>mips</code> module are connected to the testbench for monitoring
purposes.</li>
</ul>
<p>The <code>mips_tb</code> module serves as a testbench to simulate and
verify the functionality of the MIPS processor.</p>
<p>It provides the necessary inputs (clock and reset) and initializes
the memory and registers. The testbench can be modified to apply
different test cases through loading different binary converted assembly
files and allows one to monitor the processor’s behavior through the
7-segment display outputs.</p>
<h1 id="conclusion-1">Conclusion</h1>
<p>Throughout this project, I gained valuable experience in designing and
implementing a processor using Verilog. I further learned about the different
stages of the processor pipeline, including instruction fetch, decode,
execute, memory access, and write-back. I also gained a deeper
understanding of the MIPS instruction set architecture and how
instructions are encoded and executed.</p>
<p>Implementing the processor in Verilog allowed me to apply my
knowledge of digital design and hardware description languages. I
utilized various Verilog constructs, such as modules, always blocks, and
case statements, to model the behavior of the processor components. I
also learned about the importance of proper synchronization and timing
in hardware design.</p>
<p>To ensure the correctness of the processor implementation, I
developed test-benches in Verilog to verify the functionality of
individual components as well as the overall processor. The test-benches
allowed me to simulate the processor’s behavior and debug any issues
that arose during the development process.</p>
<p>In addition to the Verilog implementation, I also compared my experience with writing a similar single-cycle processor
in VHDL. This comparison provided insights into the
differences and similarities between the two hardware description
languages. While Verilog offers a more concise and flexible syntax, VHDL
provides stronger typing and more explicit component instantiation. Both
languages have their strengths and weaknesses, and the choice between
them often depends on the specific project requirements and personal
preference.</p>
</body>


