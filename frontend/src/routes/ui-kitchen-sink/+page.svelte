<script lang="ts">
	import Button from "$lib/components/ui/Button.svelte";
	import Badge from "$lib/components/ui/Badge.svelte";
	import Avatar from "$lib/components/ui/Avatar.svelte";
	import Input from "$lib/components/ui/Input.svelte";
	import Modal from "$lib/components/ui/Modal.svelte";
	import Dropdown from "$lib/components/ui/Dropdown.svelte";
	import Toast from "$lib/components/ui/Toast.svelte";
	import Form from "$lib/components/forms/Form.svelte";
	import DataTable from "$lib/components/data/DataTable.svelte";
	import Uploader from "$lib/components/ui/Uploader.svelte";
	import BarChart from "$lib/components/charts/BarChart.svelte";
	import Editor from "$lib/components/ui/Editor.svelte";
	import { toast } from "svelte-sonner";
	import { LogOut, User, Settings, MoreVertical, Plus } from "lucide-svelte";
	import { writable } from "svelte/store";

	// 1. Modal State
	let showModal = $state(false);

	// 2. Form Setup (Mock for demo)
	const formStore = writable({ username: "", email: "" });
	const errorStore = writable({ username: [], email: [] });
	const delayed = writable(false);
	
	const mockForm = {
		form: formStore,
		errors: errorStore,
		delayed: delayed,
		enhance: () => (({ result }: any) => {}),
		constraints: writable({}),
		tainted: writable(undefined),
		message: writable(undefined),
		allErrors: writable([]),
		posted: writable(false),
		submitting: writable(false),
		capture: () => {},
		restore: () => {},
		validate: async () => [],
		reset: () => {}
	} as any;

	// 3. Data Table Setup
	const data = [
		{ id: 1, name: "Andi Aryatno", email: "andi@kodia.id", role: "Owner" },
		{ id: 2, name: "Budi Santoso", email: "budi@example.com", role: "Admin" },
		{ id: 3, name: "Citra Lestari", email: "citra@example.com", role: "Member" },
		{ id: 4, name: "Deni Pratama", email: "deni@example.com", role: "Member" },
	];

	const columns = [
		{ accessorKey: "name", header: "Name" },
		{ accessorKey: "email", header: "Email" },
		{ 
			accessorKey: "role", 
			header: "Role",
			cell: (info: any) => info.getValue()
		}
	];

	// 4. Chart Data
	const chartData = [
		{ x: "Jan", y: 400 },
		{ x: "Feb", y: 300 },
		{ x: "Mar", y: 600 },
		{ x: "Apr", y: 800 },
		{ x: "May", y: 500 },
		{ x: "Jun", y: 900 },
	];
</script>

<Toast />

<div class="min-h-screen bg-hive py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-6xl mx-auto space-y-16">
		
		<!-- Header -->
		<header class="text-center space-y-4">
			<h1 class="text-6xl font-black tracking-tight text-gradient">Kodia UI Kitchen Sink</h1>
			<p class="text-xl text-slate-500 font-bold uppercase tracking-widest">Premium Svelte 5 Component Library</p>
		</header>

		<!-- 1. Buttons & Badges -->
		<section class="space-y-6">
			<h2 class="text-2xl font-black tracking-tight border-b pb-2">Buttons & Badges</h2>
			<div class="flex flex-wrap gap-4 items-center">
				<Button variant="premium">Premium Action</Button>
				<Button variant="primary">Primary</Button>
				<Button variant="secondary">Secondary</Button>
				<Button variant="outline">Outline</Button>
				<Button variant="ghost">Ghost</Button>
				<Button variant="danger" onclick={() => toast.error("Critical error occurred!")}>Danger Toast</Button>
				<Button loading variant="primary">Loading...</Button>
			</div>
			<div class="flex flex-wrap gap-3 items-center">
				<Badge variant="premium">Enterprise</Badge>
				<Badge variant="primary">Active</Badge>
				<Badge variant="success">Online</Badge>
				<Badge variant="warning">Pending</Badge>
				<Badge variant="danger">High Priority</Badge>
				<Badge variant="ghost">Draft</Badge>
			</div>
		</section>

		<!-- 2. Avatars & Dropdowns -->
		<section class="grid grid-cols-1 md:grid-cols-2 gap-12">
			<div class="space-y-6">
				<h2 class="text-2xl font-black tracking-tight border-b pb-2">Avatars</h2>
				<div class="flex items-end gap-6">
					<Avatar size="xl" src="https://github.com/shadcn.png" />
					<Avatar size="lg" src="https://github.com/shadcn.png" />
					<Avatar size="md" fallback="K" class="bg-primary/20 text-primary" />
					<Avatar size="sm" fallback="A" />
				</div>
			</div>

			<div class="space-y-6">
				<h2 class="text-2xl font-black tracking-tight border-b pb-2">Dropdowns & Modals</h2>
				<div class="flex gap-4">
					<Dropdown items={[
						{ label: "Profile", icon: User },
						{ label: "Settings", icon: Settings },
						{ label: "Logout", icon: LogOut, variant: "danger" }
					]}>
						<Button variant="outline">
							Account Options <MoreVertical class="w-4 h-4 ml-2" />
						</Button>
					</Dropdown>

					<Button variant="primary" onclick={() => showModal = true}>
						Open Premium Modal
					</Button>
				</div>
			</div>
		</section>

		<!-- 3. Form Builder -->
		<section class="space-y-6">
			<h2 class="text-2xl font-black tracking-tight border-b pb-2">Declarative Form Builder</h2>
			<div class="max-w-xl">
				<div class="card-premium">
					<Form form={mockForm} submitLabel="Register Account">
						<div class="space-y-4">
							<div>
								<label for="username" class="block text-sm font-semibold mb-2">Username</label>
								<Input id="username" placeholder="Enter your username" />
							</div>
							<div>
								<label for="email" class="block text-sm font-semibold mb-2">Email Address</label>
								<Input id="email" type="email" placeholder="you@kodia.id" />
							</div>
						</div>
					</Form>
				</div>
			</div>
		</section>

		<!-- 4. Advanced Data Table -->
		<section class="space-y-6">
			<div class="flex items-center justify-between">
				<h2 class="text-2xl font-black tracking-tight border-b pb-2">Advanced Data Table</h2>
				<Button size="sm" variant="secondary"><Plus class="w-4 h-4 mr-2" /> Add Member</Button>
			</div>
			<DataTable {data} {columns} />
		</section>

		<!-- 5. Uploader & Charts -->
		<section class="grid grid-cols-1 lg:grid-cols-2 gap-12">
			<div class="space-y-6">
				<h2 class="text-2xl font-black tracking-tight border-b pb-2">File Uploader</h2>
				<Uploader label="Upload Assets" multiple maxSize={20} />
			</div>
			<div class="space-y-6">
				<h2 class="text-2xl font-black tracking-tight border-b pb-2">Data Visualization</h2>
				<BarChart data={chartData} title="Monthly Revenue" />
			</div>
		</section>

		<!-- 6. Rich Text Editor -->
		<section class="space-y-6 pb-24">
			<h2 class="text-2xl font-black tracking-tight border-b pb-2">Rich Text Editor (WYSIWYG)</h2>
			<Editor label="Document Content" content="<p>Welcome to <strong>Kodia UI</strong>. This editor is elite and innovative.</p>" />
		</section>

	</div>
</div>

<Modal bind:open={showModal} title="Premium Deployment" description="Are you sure you want to deploy the Kodia UI library to production?">
	<div class="space-y-4">
		<p class="text-slate-600 dark:text-slate-400">
			This will enable high-fidelity components across your entire application. This action is innovative and elite.
		</p>
		<div class="bg-primary/5 p-4 rounded-kodia border border-primary/10">
			<p class="text-xs font-bold text-primary uppercase tracking-widest">System Check: Optimal</p>
		</div>
	</div>
	{#snippet footer()}
		<Button variant="ghost" onclick={() => showModal = false}>Cancel</Button>
		<Button variant="premium" onclick={() => {
			showModal = false;
			toast.success("UI Library deployed successfully!");
		}}>Confirm Deployment</Button>
	{/snippet}
</Modal>
