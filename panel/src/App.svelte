<script>
	let uploadedImage = ""
	let resultImage = ""
	let inputValue = ""
	let width = ""
	let filter = ""
	let value = ""
	//https://i.pinimg.com/736x/e2/14/fd/e214fdbdcf97bce173ca640ff67e0f06.jpg
	const loadImageFromUrl = () => {
		uploadedImage = inputValue;
		resultImage = `https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=480&height=600&url=${inputValue}&filter=${filter}${value ? `&value=${value}` :''}`
	}
	const uploadImage = (evt) => {
		const file = evt.target.files[0]
		const reader = new FileReader();
		reader.addEventListener('load',(ev)=>{
			uploadedImage = reader.result
		})
		reader.readAsDataURL(file)
	}
</script>

<main>
	<div class="container mx-auto">
		<div class="grid grid-cols-1 gap-4 md:grid-cols-3 lg:grid-cols-3">
	
			<div class="uploader-cont">
				<h1 class="text-center font-semibold text-xl my-2.5">Input image</h1>

				<div>Paste image link</div>
				<input bind:value={inputValue} placeholder="paste the image url">
				<button on:click={loadImageFromUrl}>Get image</button>
				<div class="my-2.5"> or upload an image</div>
				<img src={uploadedImage} class='image-cont' id='uploaded-image'>
				<input type="file" on:change={uploadImage}>
			</div>
			<div class="options-cont">
				<h1 class="text-center font-semibold text-xl my-2.5	">Filters</h1>

				<div class="grid grid-cols-2 gap-2 my-2.5">
					<input bind:value={width} placeholder="width">
					<input bind:value={width} placeholder="height">
					<label class="">
						Grayscale
						<input type=radio bind:group={filter} name="filter" value="grayscale">
					</label>
					<label class="">
						Sepia
						<input type=radio bind:group={filter} name="filter" value="sepia">
					</label>
					<label class="">
						Blur
						<input type=radio bind:group={filter} name="filter" value="blur">
						{#if filter === "blur"} 
							<input type="text" bind:value={value} placeholder="value">
						{/if}
					</label>
					<label class="">
						Brightness
						<input type=radio bind:group={filter} name="filter" value="brightness">
						{#if filter === "brightness"} 
							<input type="text" bind:value={value} placeholder="value">
						{/if}
					</label>
					<label class="">
						Negative
						<input type=radio bind:group={filter} name="filter" value="negative">
					</label>
					<label class="">
						Positive
						<input type=radio bind:group={filter} name="filter" value="positive">
					</label>
					<label>
						Black and white
						<input type=radio bind:group={filter} name="filter" value="blackAndWhite">
					</label>
				</div>
			</div>
			<div class="result-cont">
				<h1 class="text-center font-semibold text-xl">Output image</h1>
				<img src={resultImage} class="result-image image-cont  my-5	">
			</div>
	
		</div>
	</div>
</main>

<style>
	.image-cont {
		max-width:400px
	}
</style>