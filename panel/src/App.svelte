<script>
import { onMount } from "svelte";


	let uploadedImage = "";
	let resultImage = "";
	let inputValue = "https://www.thecanadianbazaar.com/wp-content/uploads/2017/11/Sonam-Bajwa-hot-image.jpg";
	let width = "";
	let height = "";
	let filter = "";
	let value = "";
	let widthError = false
	let heightError = false
	let blurError = false
	let hasResult = false
	let naturalWidth = ""
	let naturalHeight = ""
	let uploadimageElement
	let resizeBtn
	//https://www.thecanadianbazaar.com/wp-content/uploads/2017/11/Sonam-Bajwa-hot-image.jpg
	const loadImageFromUrl = (ev) => {
		if(!inputValue) return
		if(width <= 0  || !width) {
			widthError = true
			return
		}
		if(height <=0  || !width) {
			heightError = true
			return
		}
		if(value && filter == "blur" && value <= 0 ) {
			blurError = true
			return
		}
		ev.target.classList.add('is-loading')
		heightError = false
		widthError = false
		blurError = false
		uploadedImage = inputValue;
		resultImage = `https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=${width}&height=${height}&url=${inputValue}&filter=${filter}${
			value ? `&value=${value}` : ""
		}`;
		hasResult = true
	};

	const imageResized = (ev) => {
		resizeBtn.classList.remove('is-loading')
	}

	onMount(()=>{
		uploadimageElement.addEventListener('load',(ev)=>{
			naturalWidth = ev.target.naturalWidth
			naturalHeight = ev.target.naturalHeight
		})
		resultImgElement.addEventListener('load',(ev)=>{
			resizeBtn.classList.remove('is-loading')
		})
	})

	
</script>
<div class="app">
	<main>
		<div class="columns">
			<div class="column  uploader-cont">
				<h1 class="has-text-centered is-size-5 has-text-weight-medium">
					Input image
				</h1>
				<div class="field has-addons mt-6">
					<div class="control width-100">
					  <input bind:value={inputValue} placeholder="paste the image url" class="input" type="text" />
					</div>
					<div class="control">
					<button on:click={loadImageFromUrl} class="button is-info resize-btn" bind:this={resizeBtn}>Resize</button>
					</div>
				</div>
				<img src={uploadedImage} class="image-cont" id="uploaded-image" bind:this={uploadimageElement}/>
				{#if naturalWidth && naturalHeight} 
					<div class="mt-3 has-text-centered has-text-weight-normal">
						Natural Width : {naturalWidth}px Natural Height : {naturalHeight}px
					</div>
				{/if}
			</div>
			<div class="column options-cont">
				<h1 class="has-text-centered is-size-5 has-text-weight-medium">Image modifiers</h1>
	
				<div class="columns mt-1">
					<div class="column is-half">
						
						<div class="field">
							<label class="label">Width</label>
							<div class="control">
								<input
									bind:value={width}
									placeholder="width"
									class="input"
									type="number"
									min="1"
								/>
								{#if widthError}  
									<p class="help is-danger">Invalid width</p>
								{/if}
							</div>
						</div>
					</div>
	
					<div class="column is-half">
						<div class="field">
							<label class="label">Height</label>
							<div class="control">
								<input
									bind:value={height}
									placeholder="height"
									class="input"
									type="number"
									min="1"
								/>
								{#if heightError}  
									<p class="help is-danger">Invalid height</p>
								{/if}
							</div>
						</div>
					</div>
				</div>
				<div class="has-text-centered is-size-5 has-text-weight-medium">Filters</div>
				<div class="control mt-3">
					<div class="columns is-multiline">
						<div class="column is-one-third">
							<label class="radio">
								Grayscale
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="grayscale"
								/>
							</label>
						</div>
						<div class="column is-one-third">
							<label class="radio">
								Sepia
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="sepia"
								/>
							</label>
						</div>
						<div class="column is-one-third">
							<label class="radio">
								Blur
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="blur"
								/>
								{#if filter === "blur"}
									<input type="number" bind:value={value} placeholder="value" class="input"  />
									{#if blurError}  
										<p class="help is-danger">Value should be greater than 0</p>
									{/if}
								{/if}
							</label>
						</div>
						<div class="column is-one-third">
							<label class="radio">
								Brightness
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="brightness"
								/>
								{#if filter === "brightness"}
								<input type="number" bind:value={value} placeholder="value" class="input"/>
								{/if}
							</label>
						</div>
						<div class="column is-one-third">
							<label class="radio">
								Negative
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="negative"
								/>
							</label>
						</div>
						<div class="column is-one-third">
							<label class="radio">
								Positive
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="positive"
								/>
							</label>
						</div>
						<div class="column">
							<label class="radio">
								Black and white
								<input
									type="radio"
									bind:group={filter}
									name="filter"
									value="blackAndWhite"
								/>
							</label>
						</div>
					</div>	
				</div>
			</div>
			<div class="column result-cont">
				<h1 class="has-text-centered is-size-5 has-text-weight-medium">Output image</h1>
				<img src={resultImage} class="result-image image-cont mt-6" on:load={imageResized} />
				{#if hasResult}
				<div class="mt-2 has-text-centered">Width: {width}px Height: {height}px</div>
				<div class="mt-2 has-text-centered">
					Image Url
					<div class="mt-2">
	
						<code><a href={resultImage} target="_blank">{resultImage}</a></code>
					</div>
				</div>
				{/if}
			</div>
		</div>
	</main>
	
	<footer class="footer has-background-black">
		<div class="content has-text-centered">
		  <p class="has-text-white">
			Made with GO,SVELTE & BULMA
		  </p>
		</div>
	</footer>
</div>

<style>

	.image-cont {
		max-width: 300px;
		border: 1px solid black;
		display: block;
		margin: 0 auto;
	}
	.width-100 {
		width: 100%;
	}
	.app {
		display: flex;
		min-height: 100vh;
		flex-direction: column;
	}
	main {
		flex: 1;
		padding: 16px 8px;
	}
	footer {
		padding: 20px;
	}
</style>
