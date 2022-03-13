<script>
import { onMount } from "svelte";

	let uploadedImage = "";
	let resultImage = "";
	let inputValue = "https://cdn.myanimelist.net/r/360x360/images/anime/3/40451.jpg?s=47c23f5445fc6690845a5e69660cd8c6";
	let width = 400;
	let height = 400;
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
	let filtersSwitch = ['grayscale','sepia','negative']
	let bannerImage = "https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=200&height=200&filter="+filtersSwitch[0]+"&url=https://cdn.pixabay.com/photo/2022/01/15/02/07/windows-6938478_960_720.jpg"

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

		// let counter = 0
		// setInterval(()=>{
		// 	if (counter >= filtersSwitch.length) {
		// 		counter = 0
		// 	}
		// 	bannerImage = "https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&filter="+filtersSwitch[counter]+"&url="+inputValue
		// 	counter++
		// },8000)
	})

	
</script>
<div class="app">
	<section class="hero is-small is-black">
		
		<div class="hero-body">
			<h2 class="is-size-1 has-text-weight-bold has-text-centered">DEMO UI</h2>
			<p class="is-size-4">Refer to the <a href="https://github.com/shikhar0507/image-resizer" target="_" class="has-text-info">Github</a> README to know how to use this directly as a URL in your app</p>
			<p class="title pt-5">
				Image Manipulation (resize & filters)
			</p>
			<div class="columns">
				<div class="column is-two-third">
					<div class="is-size-4 has-text-weight-semibold">How to use</div>
					<div class="mt-2">
						<div class="is-size-5">
							Add <code>width</code> <code>height</code> & <code>url</code> in URL parameters
						</div>
						<p class="mt-4 is-size-5">Example : </p> 
						<p class="subtitle mt-2">
							<code class="is-size-6">
							<a href="https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=200&height=200&filter=grayscale&url=https://cdn.pixabay.com/photo/2022/01/15/02/07/windows-6938478_960_720.jpg" target="_">https://d3078njhubik3z.cloudfront.net/staging/imageResizer<span class="is-size-6 has-text-info">?width=300&height=300</span><span class="is-size-6 has-text-black">&url=https://cdn.pixabay.com/photo/2022/01/15/02/07/windows-6938478_960_720.jpg</span></a>
							</code>
						</p>
						<p class="subtitle is-size-6 has-text-right mt-6">
							Change params to get Image in different size and form below
						</p>						
					</div>
				</div>
				<div class="column is-one-third">
					<div class="banner-image">
						<img src={bannerImage}>
					</div>
					
				</div>
			</div>
		
		</div>
	  </section>
	<main>
		<div class="mt-2 app-cont">
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
		
					<div class="columns mt-1 is-mobile">
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
						<div class="columns is-multiline is-mobile">
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
									<div class="mt-2 is-size-6 has-text-info mb-2">(use "-" value for a darker image)</div>
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
	.app-cont {
		min-height: 500px;
	}
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
	.banner-code {
		/* padding: 20px; */
		max-width: 600px;
	}
</style>
