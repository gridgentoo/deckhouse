function updateNode(selector, storageItemName) {
  storageValue = sessionStorage.getItem(storageItemName);
  if ( storageValue && storageValue.length > 0 ) {
      $(selector).val(storageValue);
  }
}

function restoreData() {
  updateNode('#clusterdomain','dhctl-domain');
  updateNode('#resourceprefix','dhctl-prefix');
  updateNode('#sshkey','dhctl-sshkey');
  {% for preset in site.data.getting_started.data.presets %}
    {%- if preset[1].recommended %}
        sessionStorage.setItem('dhctl-preset', '{{ preset[0] }}' );
    {%- endif %}
  {%- endfor %}
{% for layout in site.data.getting_started.data.installTypes[page.platform_code].layouts %}
    {%- if layout[1].recommended %}
        sessionStorage.setItem('dhctl-layout', '{{ layout[1].name }}' );
    {%- endif %}
  {%- endfor %}
}

$(document).ready(function(){
  restoreData();
	$('#clusterdomain').change(function(){
		sessionStorage.setItem('dhctl-domain', $(this).val());
	});
	$('#resourceprefix').change(function(){
		sessionStorage.setItem('dhctl-prefix', $(this).val());
	});
	$('#sshkey').change(function(){
		sessionStorage.setItem('dhctl-sshkey', $(this).val());
	});
});